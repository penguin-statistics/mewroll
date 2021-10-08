package cmd

import (
	"fmt"
	"github.com/penguin-statistics/mewroll/internal/pkg/httpclient"
	"github.com/penguin-statistics/mewroll/internal/pkg/randomdrawer"
)

type CommentDrawerConfig struct {
	ThoughtID string
	Deduplication string
	Count int
}

type CommentDrawer struct {
	conf *CommentDrawerConfig

	client *httpclient.Client
}

func NewCommentDrawer(conf *CommentDrawerConfig) *CommentDrawer {
	return &CommentDrawer{conf: conf, client: httpclient.NewDefaultClient()}
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func (c *CommentDrawer) Draw() error {
	fmt.Printf("正在获取想法 %s 的评论内容...\n", c.conf.ThoughtID)
	comments, err := c.client.GetComments(c.conf.ThoughtID)
	if err != nil {
		return err
	}

	fmt.Printf("想法评论内容获取完毕，获取到 %d 条未经处理的评论数量。开始进行预处理...\n", len(comments.Entries))

	newEntries := []httpclient.MewThoughtComment{}
	for _, entry := range comments.Entries {
		if !entry.Deleted {
			newEntries = append(newEntries, entry)
		}
	}
	comments.Entries = newEntries

	fmt.Printf("预处理完毕，获取到 %d 条未经去重的评论数量。开始进行去重...\n", len(comments.Entries))

	var candidatePostIndexes []int

	authorHasPostIndexes := map[string][]int{}
	for _, entry := range comments.Entries {
		//if entry.Deleted {
		//	continue
		//}

		if _, ok := authorHasPostIndexes[entry.AuthorId]; ok {
			authorHasPostIndexes[entry.AuthorId] = append(authorHasPostIndexes[entry.AuthorId], entry.Index)
		} else {
			authorHasPostIndexes[entry.AuthorId] = []int{entry.Index}
		}
	}

	if c.conf.Deduplication == "eliminate" {
		fmt.Println("正在使用 eliminate 模式去重...")
		for authorId, postIndexes := range authorHasPostIndexes {
			if len(postIndexes) == 1 {
				candidatePostIndexes = append(candidatePostIndexes, postIndexes[0])
			} else {
				fmt.Printf("  - 用户 %s 由于发送了多于 1 条的想法数量（发送楼层：%v），已被从奖池内移除\n", authorId, postIndexes)
			}
		}
	} else if c.conf.Deduplication == "single" {
		fmt.Println("正在使用 single 模式去重...")
		for _, postIndexes := range authorHasPostIndexes {
			candidatePostIndexes = append(candidatePostIndexes, postIndexes[0])
		}
	} else {
		for _, postIndexes := range authorHasPostIndexes {
			candidatePostIndexes = append(candidatePostIndexes, postIndexes...)
		}
	}

	total := len(candidatePostIndexes)

	fmt.Printf("想法评论内容去重完毕，去重后的评论数量为 %d。开始抽取...\n", total)

	for i := 0; i < c.conf.Count; i++ {
		currentTotal := len(candidatePostIndexes)

		//fmt.Printf("剩余 %d 条评论。开始抽取...\n", currentTotal)

		drawResult := randomdrawer.Draw(0, currentTotal)

		var comment httpclient.MewThoughtComment
		for _, entry := range comments.Entries {
			if entry.Index == candidatePostIndexes[drawResult] {
				comment = entry
			}
		}

		user := comments.Objects.Users[comment.AuthorId]

		fmt.Printf("  - #%d：由 %s (@%s, %s) 于楼层 #%d 发表的内容「%s」\n", i, user.Name, user.Username, user.Id, comment.Index, comment.Content)

		candidatePostIndexes = remove(candidatePostIndexes, drawResult)
	}

	return nil
}
