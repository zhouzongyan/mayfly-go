package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"mayfly-go/internal/db/config"
	"mayfly-go/internal/db/domain/entity"
	"mayfly-go/internal/db/domain/repository"
	"mayfly-go/pkg/base"
	"mayfly-go/pkg/model"
	"os"
	"path/filepath"
)

type DbTransferFile interface {
	base.App[*entity.DbTransferFile]

	// GetPageList 分页获取数据库实例
	GetPageList(condition *entity.DbTransferFileQuery, pageParam *model.PageParam, toEntity any, orderBy ...string) (*model.PageResult[any], error)

	Save(ctx context.Context, instanceEntity *entity.DbTransferFile) error

	Delete(ctx context.Context, id ...uint64) error

	GetFilePath(ent *entity.DbTransferFile) string
}

var _ DbTransferFile = (*dbTransferFileAppImpl)(nil)

type dbTransferFileAppImpl struct {
	base.AppImpl[*entity.DbTransferFile, repository.DbTransferFile]
}

func (app *dbTransferFileAppImpl) InjectDbTransferFileRepo(repo repository.DbTransferFile) {
	app.Repo = repo
}

func (app *dbTransferFileAppImpl) GetPageList(condition *entity.DbTransferFileQuery, pageParam *model.PageParam, toEntity any, orderBy ...string) (*model.PageResult[any], error) {
	return app.GetRepo().GetPageList(condition, pageParam, toEntity, orderBy...)
}

func (app *dbTransferFileAppImpl) Save(ctx context.Context, taskEntity *entity.DbTransferFile) error {
	var err error
	if taskEntity.Id == 0 {
		err = app.Insert(ctx, taskEntity)
	} else {
		err = app.UpdateById(ctx, taskEntity)
	}
	return err
}

func (app *dbTransferFileAppImpl) Delete(ctx context.Context, id ...uint64) error {

	arr, err := app.GetByIds(id, "task_id", "file_uuid")
	if err != nil {
		return err
	}

	// 删除对应的文件
	for _, file := range arr {
		_ = os.Remove(app.GetFilePath(file))
	}

	// 删除数据
	return app.DeleteById(ctx, id...)
}

func (app *dbTransferFileAppImpl) GetFilePath(ent *entity.DbTransferFile) string {
	brc := config.GetDbBackupRestore()
	if ent.FileUuid == "" {
		ent.FileUuid = uuid.New().String()
	}

	filePath := filepath.Join(fmt.Sprintf("%s/%d/%s.sql", brc.TransferPath, ent.TaskId, ent.FileUuid))

	return filePath
}