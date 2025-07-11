package selectallacrosspages

import (
	"context"
	"example.com/m/util/collection"
	"example.com/m/util/paginator"
	"example.com/m/util/utilerror"
)

type SelectAllAcrossPagesTemplate struct {
	implement SelectAllAcrossPagesInterface
}

func NewSelectAllAcrossPagesTemplate(implement SelectAllAcrossPagesInterface) *SelectAllAcrossPagesTemplate {
	return &SelectAllAcrossPagesTemplate{
		implement: implement,
	}
}

type SelectAllAcrossPagesTemplateRequest struct {
	Condition            interface{} //搜索栏选项
	IsAllChoose          bool        //是否全选
	LastTimeChosenIdList []string    //用户之前选择的所有id
	ThisTimeChosenIdList []string    //本次的勾选
	NoChosenIdList       []string    //本次的反选
	SortIdParam          interface{} //id的排序规则
	FillFieldParam       interface{} //填充列表页相关字段需要用到的参数
	PageNo               int64
	Count                int64
}

type SelectAllAcrossPagesTemplateResponse struct {
	ChosenIdList         []string    //用户勾选的所有id
	ThisPageList         interface{} //列表展示的数据
	ThisPageChosenIdList []string    //本页需要勾选的数据
	Total                int64
	PageNo               int64
	Count                int64
}

func (t *SelectAllAcrossPagesTemplate) Select(ctx context.Context, req *SelectAllAcrossPagesTemplateRequest, includeDeleted bool) (*SelectAllAcrossPagesTemplateResponse, *utilerror.UtilError) {
	allIdList, selectErr := t.implement.SelectIdByCondition(ctx, req.Condition, includeDeleted)
	if selectErr != nil {
		return nil, selectErr.Mark()
	}
	chosenIdSet := collection.NewStringSet(req.LastTimeChosenIdList...)
	if req.IsAllChoose {
		chosenIdSet.Add(allIdList...)
	} else {
		chosenIdSet.Add(req.ThisTimeChosenIdList...)
	}
	if len(req.NoChosenIdList) > 0 {
		chosenIdSet.Remove(req.NoChosenIdList...)
	}

	total := len(allIdList)
	allIdList, sortErr := t.implement.SortId(ctx, allIdList, req.SortIdParam)
	if sortErr != nil {
		return nil, sortErr.Mark()
	}
	paginator.PageList(&allIdList, req.PageNo, req.Count)
	thisPageChosenIdSet := collection.NewStringSet(allIdList...)
	thisPageChosenIdSet = thisPageChosenIdSet.InterSet(chosenIdSet)
	response := &SelectAllAcrossPagesTemplateResponse{
		ChosenIdList:         chosenIdSet.ToSlice(),
		Total:                int64(total),
		PageNo:               req.PageNo,
		Count:                req.Count,
		ThisPageChosenIdList: thisPageChosenIdSet.ToSlice(),
	}

	thisPageList, fillErr := t.implement.FillFields(ctx, allIdList, req.FillFieldParam, includeDeleted)
	if fillErr != nil {
		return nil, fillErr.Mark()
	}
	response.ThisPageList = thisPageList
	return response, nil
}

type SelectAllAcrossPagesInterface interface {
	// SelectIdByCondition 根据condition查出满足条件的所有数据的id
	SelectIdByCondition(ctx context.Context, condition interface{}, includeDeleted bool) ([]string, *utilerror.UtilError)
	// SortId id的排序规则
	SortId(ctx context.Context, allIdList []string, sortParam interface{}) ([]string, *utilerror.UtilError)
	// FillFields 根据thisPageIdList中的id查出列表页的数据，记得排序
	FillFields(ctx context.Context, thisPageIdList []string, fillFieldParam interface{}, includeDeleted bool) (interface{}, *utilerror.UtilError)
}
