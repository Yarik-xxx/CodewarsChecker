package CodeWars

// RespUser Request response structure https://www.codewars.com/api/v1/users/{user}
type RespUser struct {
	Username            string         `json:"username"`
	Name                string         `json:"name"`
	Honor               int            `json:"honor"`
	LeaderboardPosition int            `json:"leaderboardPosition"`
	Ranks               Ranks          `json:"ranks"`
	CodeChallenges      CodeChallenges `json:"codeChallenges"`
	Reason              string         `json:"reason"`
}

type Ranks struct {
	Overall   Rank            `json:"overall"`
	Languages map[string]Rank `json:"languages"`
}

type Rank struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type CodeChallenges struct {
	TotalCompleted int `json:"totalCompleted"`
}

// RespCompletedChallenges Request response structure https://www.codewars.com/api/v1/users/{user}/code-challenges/completed?page={page}
type RespCompletedChallenges struct {
	TotalPages     int              `json:"totalPages"`
	TotalItems     int              `json:"totalItems"`
	ListChallenges []ListChallenges `json:"data"`
	Reason         string           `json:"reason"`
}

type ListChallenges struct {
	ID string `json:"id"`
	//CompletedAt        time.Time `json:"completedAt"` // In future
	CompletedLanguages []string `json:"completedLanguages"`
}

type respGetListChallenges struct {
	challenges RespCompletedChallenges
	err        error
}

type RespCodeChallenge struct {
	ID string `json:"id"`
	//Category    string   `json:"category"`  // In future
	Rank struct {
		Name string `json:"name"`
	} `json:"rank"`
	Reason string `json:"reason"`
}
