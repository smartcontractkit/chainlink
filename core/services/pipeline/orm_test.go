package pipeline

import (
	"context"
	"fmt"
	"testing"
	"time"

	gormv1 "github.com/jinzhu/gorm"
	p "github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestClaimTaskScale(t *testing.T) {
	dsn := "host=localhost user=postgres password=node dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	dbname := "load"
	err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname)).Error
	require.NoError(t, err)
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error
	require.NoError(t, err)
	dbg, err := gormv1.Open("postgres", "postgres://postgres:node@localhost:5432/load?sslmode=disable")
	require.NoError(t, err)
	err = migrations.GORMMigrate(dbg)
	require.NoError(t, err)
	dsn = "host=localhost user=postgres password=node dbname=load port=5432 sslmode=disable"
	db2, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		// Create pipeline specs
		sp := Spec{DotDagSource: fmt.Sprintf("test%d", i)}
		err := db2.Create(&sp).Error
		require.NoError(t, err)
		// Create pipeline runs
		pr := Run{
			PipelineSpecID: sp.ID,
		}
		if i < 4 {
			n := time.Now()
			pr.FinishedAt = &n
			pr.Outputs = &JSONSerializable{Val: 10}
			pr.Errors = &JSONSerializable{}
		}
		err = db2.Create(&pr).Error
		require.NoError(t, err)
		t.Log("created pipeline run", pr.ID)
		// Create pipeline task specs
		var prev int32
		var preds = make(map[int32]int32)
		var spids []int32
		// Insert them backwards
		for k := 0; k < 3; k++ {
			ts := TaskSpec{
				PipelineSpecID: sp.ID,
				Type:           "",
			}
			if prev != 0 {
				ts.SuccessorID = null.NewInt(int64(prev), true)
			} else {
				ts.SuccessorID = null.NewInt(0, false)
			}
			err = db2.Create(&ts).Error
			if prev != 0 {
				preds[prev] = ts.ID
			}
			prev = ts.ID
			spids = append([]int32{ts.ID}, spids...)
		}
		fmt.Println(spids)
		// For each job run, create the task runs
		// including the predessor task run info
		// That would mean we need to create the task runs inorder of successors
		for j := 0; j < 1000; j++ {
			var last int64
			for k := 0; k < 3; k++ {
				tr := TaskRun{
					ID:                 0,
					PipelineRunID:      pr.ID,
					Output:             nil,
					Error:              null.String{},
					PipelineTaskSpecID: spids[k],
				}
				if last != 0 {
					tr.PredecessorTaskRunIds = []int64{last}
				}
				if k < 2 {
					n := time.Now()
					tr.FinishedAt = &n
				}
				require.NoError(t, db2.Create(&tr).Error)
				last = tr.ID
			}
			last = 0
		}
	}

	o := orm{db: dbg}

	var ptRun TaskRun
	s := time.Now()
	err = p.GormTransaction(context.Background(), o.db, func(tx *gormv1.DB) error {
		return tx.Raw(`
				select * from pipeline_task_runs where id in (
					select runs.id from pipeline_task_runs as runs
					left join pipeline_task_runs as preds on preds.id = any (runs.predecessor_task_run_ids)
					where runs.finished_at is null
					group by runs.id
					having (bool_and(preds.finished_at is not null) or count(preds) = 0)
				)
				limit 1
				FOR UPDATE OF pipeline_task_runs SKIP LOCKED;`).Scan(&ptRun).Error
	})
	require.NoError(t, err)
	e := time.Now()
	t.Log("runtime", e.Sub(s))
	t.Log(ptRun.ID, ptRun.PipelineRunID)

	s = time.Now()
	err = p.GormTransaction(context.Background(), o.db, func(tx *gormv1.DB) error {
		return tx.Raw(`
		   SELECT * from pipeline_task_runs WHERE id IN (
		       SELECT pipeline_task_runs.id FROM pipeline_task_runs
		           INNER JOIN pipeline_task_specs ON pipeline_task_runs.pipeline_task_spec_id = pipeline_task_specs.id
		           LEFT JOIN pipeline_task_specs AS predecessor_specs ON predecessor_specs.successor_id = pipeline_task_specs.id
		           LEFT JOIN pipeline_task_runs AS predecessor_unfinished_runs
							ON (predecessor_specs.id = predecessor_unfinished_runs.pipeline_task_spec_id
		                 		AND pipeline_task_runs.pipeline_run_id = predecessor_unfinished_runs.pipeline_run_id)
		       WHERE pipeline_task_runs.finished_at IS NULL
		       GROUP BY (pipeline_task_runs.id)
		       HAVING (
		           bool_and(predecessor_unfinished_runs.finished_at IS NOT NULL)
		           OR
		           count(predecessor_unfinished_runs.id) = 0
		       )
		   )
		   LIMIT 1
			FOR UPDATE OF pipeline_task_runs SKIP LOCKED;
		`).Scan(&ptRun).Error
	})
	require.NoError(t, err)
	e = time.Now()
	t.Log("runtime2", e.Sub(s))
	t.Log(ptRun.ID, ptRun.PipelineRunID)
}
