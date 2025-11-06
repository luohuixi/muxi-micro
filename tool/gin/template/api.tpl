type Class {
    Id int64 `json:"id" binding:"require"` //课程ID
    Name string `json:"name"` //课程名称
    Teacher string `json:"teacher"` //授课教师
    Credit int `json:"credit"` //学分
    Description string `json:"description"` //课程简介
}

type (
    CreateClassReq {
        Name string `json:"name"` //课程名称
        Teacher string `json:"teacher"` //授课教师
        Credit int `json:"credit"` //学分
        Description string `json:"description"` //课程简介
    }

    CreateClassResp {
        ClassId int64 `json:"class_id"` //新建课程的ID
        Message string `json:"message"` //提示信息
    }
)

type (
    GetClassReq {
        Id int64 `form:"id"` //课程ID
    }

    GetClassResp {
        Class Class `json:"class"` //课程详情
    }
)

@server (
    prefix: /v1
    group: class
)

service class {
    @doc (
        summary: "创建课程"
        description: "管理员创建新的课程信息。"
        tag: "class"
        accept: "json"
        produce: "json"
        success: "200 {object} CreateClassResp"
    )
    @handler CreateClass
    post /class (CreateClassReq) returns (CreateClassResp)

    @doc (
        summary: "获取课程详情"
        description: "根据课程ID获取课程的详细信息。"
    )
    @handler GetClass
    get /class/:id returns (GetClassResp)
}
