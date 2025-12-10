## 1. 目标
为考勤历史页面添加员工姓名搜索功能，提高系统查询灵活性和用户体验。

## 2. 实施
- [ ] 2.1 前端： `views/attendance_history_manage.html` ，在搜索表单中添加员工姓名输入框，修改JavaScript查询逻辑以支持姓名和工号的组合搜索。
- [ ] 2.2 后端服务： `service/attend_record.go` ，修改`GetAttendRecordHistoryByStaffId`函数，添加按员工姓名查询的SQL逻辑，支持姓名和工号的组合查询。
- [ ] 2.3 后端处理： `handler/attend.go` ，可能需要修改`GetAttendRecordHistoryByStaffId`处理函数，以接收和处理员工姓名参数。
- [ ] 2.4 路由配置： `main.go` ，如需要，添加新的路由支持员工姓名搜索。