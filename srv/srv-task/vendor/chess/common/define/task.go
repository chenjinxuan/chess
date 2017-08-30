package define

const (
	TaskLoopHandleGameOverRedisKey    = "task_loop_handle_game_over_redis_key"
	TaskLoopHandlePlayerEventRedisKey = "task_loop_handle_player_event_redis_key"
        TaskUpsetRedisKey                 = "task_upset_redis_key"
	TodayTask                         = 1
	WeekTask                          = 2
	PermanentTask                     = 3
        RequiredCommon                    = 0 //只需要参加
        RequiredWin                       = 1 //要求赢
        RequiredTotalBalance              = 2  //要求金币余额

)
