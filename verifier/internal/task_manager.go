package internal

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slog"

	"github.com/game-of-nfts/gon-toolbox/verifier/internal/verifier"
)

type (
	Options struct {
		TaskIDs []string
	}

	Task struct {
		TaskNo string
		Params any
		vf     verifier.Verifier
	}

	TaskManager struct {
		tasks     []Task
		user      verifier.UserInfo
		wg        *sync.WaitGroup
		outputDir string
		resultCh  chan *verifier.Respone
		stopCh    chan int
	}
)

func NewTaskManager(evidenceFile string, opts *Options) (*TaskManager, error) {
	tm := &TaskManager{
		wg:       &sync.WaitGroup{},
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

	go tm.receive()
	for _, task := range tm.tasks {
		tm.wg.Add(1)
		go func(task Task) {
			defer tm.wg.Done()
			task.vf.Do(*&verifier.Request{
				TaskNo: task.TaskNo,
				User:   tm.user,
				Params: task.Params,
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
		case <-tm.stopCh:
			fileName := filepath.Join(tm.outputDir, fmt.Sprintf("%s.xlsx", tm.user.TeamName))
			if err := f.SaveAs(fileName); err != nil {
				slog.Error("Save file error", err)
			}
		}
	}
}

func (tm *TaskManager) stop() {
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

	tm.outputDir = filepath.Dir(evidenceFile)
	return tm.buildVerifier(evidence, opts)
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
	tm.user = verifier.UserInfo{
		TeamName: columns[0],
		Github:   "", //TODO
		Address: map[string]string{
			"i": columns[1],
			"s": columns[2],
			"j": columns[3],
			"u": columns[4],
			"o": columns[5],
		},
	}
	return nil
}

func (tm *TaskManager) buildVerifier(evidence *excelize.File, opts *Options) error {
	//TODO
	evidence.GetSheetList()
	return nil
}
