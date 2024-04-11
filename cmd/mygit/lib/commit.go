package lib

import "fmt"

func prepareAuthor(name, email string) (string, string) {
	if name == "" {
		name = "mygit"
	}
	if email == "" {
		email = "mygit@example.com"
	}
	return name, email
}

func CreateCommit(treeHash, parentHash []byte, message, authorName, authorEmail string) []byte {
	authorName, authorEmail = prepareAuthor(authorName, authorEmail)

	commitContent := fmt.Sprintf("tree %x\n", treeHash)
	if parentHash != nil {
		commitContent += fmt.Sprintf("parent %x\n", parentHash)
	}
	commitContent += fmt.Sprintf("author %s <%s> 1620000000 +0000\n", authorName, authorEmail)
	commitContent += fmt.Sprintf("committer %s <%s> 1620000000 +0000\n", authorName, authorEmail)
	commitContent += fmt.Sprintf("\n%s\n", message)

	commitHeader := fmt.Sprintf("commit %d", len(commitContent)+1)

	return append(append([]byte(commitHeader), 0), []byte(commitContent)...)
}
