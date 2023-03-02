package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slog"

	"github.com/game-of-nfts/gon-toolbox/verifier/internal/chain"
	"github.com/game-of-nfts/gon-toolbox/verifier/internal/verifier"
)

type (
	Options struct {
		TaskNos []string
	}

	Task struct {
		taskNo string
		params any
		vf     verifier.Verifier
	}

	TaskManager struct {
		tasks    []Task
		user     verifier.UserInfo
		vr       *verifier.Registry
		wg       *sync.WaitGroup
		baseDir  string
		resultCh chan *verifier.Respone
		stopCh   chan int
	}
)

func NewTaskManager(evidenceFile string, opts *Options) (*TaskManager, error) {
	tm := &TaskManager{
		wg:       &sync.WaitGroup{},
		vr:       verifier.NewRegistry(chain.NewRegistry()),
		resultCh: make(chan *verifier.Respone, 10),
		stopCh:   make(chan int),
	}

	if err := tm.loadEvidence(evidenceFile, opts); err != nil {
		return nil, err
	}
	return tm, nil
}

func (tm *TaskManager) Process() {
	if len(tm.tasks) == 0 {
		slog.Info("no task process")
		return
	}
	slog.Info("start to verify", "TeamName", tm.user.TeamName)
	go tm.receive()
	for _, task := range tm.tasks {
		tm.wg.Add(1)
		go func(task Task) {
			defer tm.wg.Done()
			slog.Info("verify rule", "TeamName", tm.user.TeamName, "TaskNo", task.taskNo)
			task.vf.Do(*&verifier.Request{
				TaskNo: task.taskNo,
				User:   tm.user,
				Params: task.params,
			}, tm.resultCh)
		}(task)
	}
	tm.wg.Wait()
	tm.stop()
	return
}

func (tm *TaskManager) receive() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("close file error", err)
		}
	}()

	sheetName := "result"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		slog.Error("NewSheet error", err)
		return
	}
	f.SetActiveSheet(index)

	var rowIdx int
	for {
		select {
		case result := <-tm.resultCh:
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), result.TaskNo)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), result.TeamName)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), result.Point)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), result.Reason)
			rowIdx++
		case <-tm.stopCh:
			fileName := filepath.Join(tm.baseDir, fmt.Sprintf("%s.xlsx", tm.user.Github))
			if err := f.SaveAs(fileName); err != nil {
				slog.Error("Save file error", err)
			}
		}
	}
}

func (tm *TaskManager) stop() {
	slog.Info("verify finish", "TeamName", tm.user.TeamName)
	tm.stopCh <- 1
}

func (tm *TaskManager) loadEvidence(evidenceFile string, opts *Options) error {
	evidence, err := excelize.OpenFile(evidenceFile)
	if err != nil {
		return err
	}
	if err := tm.loadUserInfo(evidence); err != nil {
		return err
	}

	tm.baseDir = filepath.Dir(evidenceFile)
	return tm.buildTask(evidence, opts)
}

func (tm *TaskManager) loadUserInfo(evidence *excelize.File) error {
	rows, err := evidence.GetRows("Info")
	if err != nil {
		return err
	}

	if len(rows) != 2 {
		return errors.New("invalid evidence template")
	}

	columns := rows[1]
	github := columns[0]
	if len(tm.baseDir) > 0 {
		paths := strings.Split(tm.baseDir, string(os.PathSeparator))
		github = paths[len(paths)-1]
	}

	tm.user = verifier.UserInfo{
		TeamName: columns[0],
		Github:   github,
		Address: map[string]string{
			chain.ChainIdAbbreviationIris:     columns[1],
			chain.ChainIdAbbreviationStars:    columns[2],
			chain.ChainIdAbbreviationJuno:     columns[3],
			chain.ChainIdAbbreviationUptick:   columns[4],
			chain.ChainIdAbbreviationOmniflix: columns[5],
		},
	}
	return nil
}

func (tm *TaskManager) buildTask(evidence *excelize.File, opts *Options) error {
	taskNos := evidence.GetSheetList()
	if len(opts.TaskNos) != 0 {
		taskNos = opts.TaskNos
	}
	for _, taskNo := range taskNos {
		rowsCols, err := evidence.GetRows(taskNo)
		if err != nil {
			return err
		}

		vf := tm.vr.Get(taskNo)
		params, err := vf.BuildParams(rowsCols[1:])
		if err != nil {
			return err
		}

		tm.tasks = append(tm.tasks, Task{
			taskNo: taskNo,
			params: params,
			vf:     vf,
		})
	}
	return nil
}
