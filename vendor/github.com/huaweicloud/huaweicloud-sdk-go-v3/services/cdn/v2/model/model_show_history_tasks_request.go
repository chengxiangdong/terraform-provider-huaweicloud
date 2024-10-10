package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowHistoryTasksRequest Request Object
type ShowHistoryTasksRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子帐号调用接口时，该参数必传。  您可以通过调用企业项目管理服务（EPS）的查询企业项目列表接口（ListEnterpriseProject）查询企业项目id。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 单页最大数量，取值范围为1-10000。page_size和page_number必须同时传值。默认值30。
	PageSize *int32 `json:"page_size,omitempty"`

	// 当前查询第几页，取值范围为1-65535。默认值1。
	PageNumber *int32 `json:"page_number,omitempty"`

	// 任务状态。 task_inprocess 表示任务处理中，task_done表示任务完成。
	Status *ShowHistoryTasksRequestStatus `json:"status,omitempty"`

	// 查询起始时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	StartDate *int64 `json:"start_date,omitempty"`

	// 查询结束时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	EndDate *int64 `json:"end_date,omitempty"`

	// 用来排序的字段，支持的字段有“task_type”：任务的类型，“total”：url总数，“processing”：处理中的url个数， “succeed”：成功处理的url个数，“failed”：处理失败的url个数，“create_time”：任务的创建时间。order_field和order_type必须同时传值，否则使用默认值\"create_time\" 和 \"desc\"：降序。
	OrderField *string `json:"order_field,omitempty"`

	// desc：降序，或者asc：升序。默认值desc。
	OrderType *string `json:"order_type,omitempty"`

	// 默认是文件file。file：文件,directory：目录。
	FileType *ShowHistoryTasksRequestFileType `json:"file_type,omitempty"`

	// 任务类型，refresh：刷新任务；preheating：预热任务
	TaskType *ShowHistoryTasksRequestTaskType `json:"task_type,omitempty"`
}

func (o ShowHistoryTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowHistoryTasksRequest struct{}"
	}

	return strings.Join([]string{"ShowHistoryTasksRequest", string(data)}, " ")
}

type ShowHistoryTasksRequestStatus struct {
	value string
}

type ShowHistoryTasksRequestStatusEnum struct {
	TASK_INPROCESS ShowHistoryTasksRequestStatus
	TASK_DONE      ShowHistoryTasksRequestStatus
}

func GetShowHistoryTasksRequestStatusEnum() ShowHistoryTasksRequestStatusEnum {
	return ShowHistoryTasksRequestStatusEnum{
		TASK_INPROCESS: ShowHistoryTasksRequestStatus{
			value: "task_inprocess",
		},
		TASK_DONE: ShowHistoryTasksRequestStatus{
			value: "task_done",
		},
	}
}

func (c ShowHistoryTasksRequestStatus) Value() string {
	return c.value
}

func (c ShowHistoryTasksRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowHistoryTasksRequestStatus) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type ShowHistoryTasksRequestFileType struct {
	value string
}

type ShowHistoryTasksRequestFileTypeEnum struct {
	FILE      ShowHistoryTasksRequestFileType
	DIRECTORY ShowHistoryTasksRequestFileType
}

func GetShowHistoryTasksRequestFileTypeEnum() ShowHistoryTasksRequestFileTypeEnum {
	return ShowHistoryTasksRequestFileTypeEnum{
		FILE: ShowHistoryTasksRequestFileType{
			value: "file",
		},
		DIRECTORY: ShowHistoryTasksRequestFileType{
			value: "directory",
		},
	}
}

func (c ShowHistoryTasksRequestFileType) Value() string {
	return c.value
}

func (c ShowHistoryTasksRequestFileType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowHistoryTasksRequestFileType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}

type ShowHistoryTasksRequestTaskType struct {
	value string
}

type ShowHistoryTasksRequestTaskTypeEnum struct {
	REFRESH    ShowHistoryTasksRequestTaskType
	PREHEATING ShowHistoryTasksRequestTaskType
}

func GetShowHistoryTasksRequestTaskTypeEnum() ShowHistoryTasksRequestTaskTypeEnum {
	return ShowHistoryTasksRequestTaskTypeEnum{
		REFRESH: ShowHistoryTasksRequestTaskType{
			value: "refresh",
		},
		PREHEATING: ShowHistoryTasksRequestTaskType{
			value: "preheating",
		},
	}
}

func (c ShowHistoryTasksRequestTaskType) Value() string {
	return c.value
}

func (c ShowHistoryTasksRequestTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowHistoryTasksRequestTaskType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
