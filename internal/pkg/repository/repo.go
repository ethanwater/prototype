package repository

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Commit struct {
	Hash string
	Date string

	Tree      string
	Parent    string
	Committer string
}

const Count int = 100

func main() {
	commits := make([]Commit, 0, 100)
	maxcount := "--max-count=" + strconv.Itoa(Count)
	stdout, err := exec.Command("git", "hist", maxcount).Output()
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(stdout), "\n")

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, entry := range lines {
		entry = removeNonAlphanumeric(entry)
		if len(entry) < 2 {
			continue
		}

		// Process commit data in parallel
		wg.Add(1)
		go func(entry string) {
			defer wg.Done()
			var commit Commit
			fields := strings.Fields(entry)
			if len(fields) >= 2 {
				hash := fields[1]
				data := commitData(hash)

				date := fields[0]
				tree := data[0]
				parent := data[1]
				var committer string
				if len(parent) < 30 {
					committer = "-----BEGIN"
				} else {
					committer = data[2]
				}

				commit = Commit{hash, date, tree, parent, committer}
				// Append to commits under a lock to avoid race conditions
				mutex.Lock()
				commits = append(commits, commit)
				mutex.Unlock()
			}
		}(entry)
	}
	wg.Wait()
	sortCommits(commits)

	blob := commitTreeData(commits[1].Tree)
	fmt.Println(commits[1])

	for _, x := range blob {
		fmt.Println(x)
	}
}

var regex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)

func removeNonAlphanumeric(input string) string {
	return regex.ReplaceAllString(input, "")
}

type ByDate []Commit

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[j].Date < a[i].Date }

func sortCommits(commits []Commit) {
	sort.Sort(ByDate(commits))
	for index, commit := range commits {
		if strings.Contains(commit.Committer, "BEGIN") {
			commits = append(commits[:index], commits[index+1:]...)
			commits = append(commits, commit)
		}
	}
}

func commitData(hash string) []string {
	var commit_data []string
	stdout, err := exec.Command("git", "cat-file", "-p", hash).Output()
	if err != nil {
		fmt.Println(err)
	}

	data := strings.Split(string(stdout), "\n")[:4]
	for _, entry := range data {
		fields := strings.Fields(strings.ReplaceAll(entry, "*", ""))
		if len(fields) > 1 {
			commit_data = append(commit_data, fields[1])
		}
	}
	return commit_data
}

type Blob struct {
	filetype   string
	binaryhash string
	filename   string
}

func commitTreeData(tree string) []Blob {
	var commit_tree_data []Blob
	stdout, err := exec.Command("git", "cat-file", "-p", tree).Output()
	if err != nil {
		fmt.Println(err)
	}

	data := strings.Split(string(stdout), "\n")
	for _, entry := range data {
		fields := strings.Fields(entry)
		if len(fields) > 1 {
			commit_tree_data = append(commit_tree_data, Blob{fields[1], fields[2], fields[3]})
		}
	}

	return commit_tree_data
}

