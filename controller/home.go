package controller

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type BranchInfo struct {
	Name       string
	LastCommit string
}

// Helper: extract owner and repo from GitHub URL
func parseGitHubURL(url string) (string, string) {
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

func getGitHubClientWithToken(token string) (*github.Client, context.Context) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	ctx := context.Background()
	client := github.NewClient(tc)
	return client, ctx
}

func Home(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	token, url := sess.Get("token"), sess.Get("url")
	if token == nil || url == nil {
		return ctx.Redirect("/login")
	}

	owner, repo := parseGitHubURL(url.(string))
	if owner == "" || repo == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid GitHub URL format")
	}

	client, c := getGitHubClientWithToken(token.(string))

	page := ctx.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	opts := &github.BranchListOptions{
		ListOptions: github.ListOptions{Page: page, PerPage: 10},
	}

	branches, resp, err := client.Repositories.ListBranches(c, owner, repo, opts)
	if err != nil {
		log.Printf("GitHub API error: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to list branches")
	}

	var branchInfos []BranchInfo
	for _, b := range branches {
		if b.Name != nil && *b.Name != "main" && *b.Name != "develop" {
			branchInfos = append(branchInfos, BranchInfo{
				Name:       b.GetName(),
				LastCommit: b.GetCommit().GetSHA(),
			})
		}
	}

	totalPages := 1
	if resp.LastPage > 0 {
		totalPages = resp.LastPage
	}

	return ctx.Render("main", fiber.Map{
		"branches":   branchInfos,
		"page":       page,
		"totalPages": totalPages,
	})
}

func DeleteBranches(ctx *fiber.Ctx) error {
	sess, err := store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Session error")
	}

	token, url := sess.Get("token"), sess.Get("url")
	if token == nil || url == nil {
		return ctx.Redirect("/login")
	}

	owner, repo := parseGitHubURL(url.(string))
	if owner == "" || repo == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid GitHub URL format")
	}

	client, c := getGitHubClientWithToken(token.(string))

	formBranches := ctx.FormValue("branches")
	if formBranches == "" {
		return ctx.Redirect("/")
	}

	branches := strings.Split(formBranches, ",")
	var failed []string

	for _, branch := range branches {
		ref := "refs/heads/" + strings.TrimSpace(branch)
		_, err := client.Git.DeleteRef(c, owner, repo, ref)
		if err != nil {
			log.Printf("Failed to delete branch %s: %v", branch, err)
			failed = append(failed, branch)
		}
	}

	if len(failed) > 0 {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Some branches could not be deleted: " + strings.Join(failed, ", "))
	}

	return ctx.Redirect("/")
}
