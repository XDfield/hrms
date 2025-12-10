## 1. 目标
为考勤历史页面添加姓名搜索功能，使用户可以通过员工姓名查询考勤历史记录。

## 2. 实施
- [ ] 2.1 前端：`views/attendance_history_manage.html` ，添加姓名搜索输入框并修改搜索逻辑。
- [ ] 2.2 后端：`handler/attend.go` ，修改`GetAttendRecordHistoryByStaffId`函数，添加支持姓名搜索的参数处理。
- [ ] 2.3 服务层：`service/attend_record.go` ，添加支持姓名搜索的查询方法。
- [ ] 2.4 路由：`main.go` ，如有需要，添加新的路由处理姓名搜索参数。