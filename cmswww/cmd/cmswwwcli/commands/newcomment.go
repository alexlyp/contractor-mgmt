package commands

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/decred/contractor-mgmt/cmswww/api/v1"
	"github.com/decred/contractor-mgmt/cmswww/cmd/cmswwwcli/config"
)

// Help message displayed for the command 'politeiawwwcli help newcomment'
var NewCommentCmdHelpMsg = `newcomment "token" "comment"

Comment on proposal as logged in user. 

Arguments:
1. token       (string, required)   Proposal censorship token
2. comment     (string, required)   Comment
3. parentID    (string, required if replying to comment)  Id of commment

Request:
{
  "token":       (string)  Censorship token
  "parentid":    (string)  Id of comment (defaults to '0' (top-level comment))
  "comment":     (string)  Comment
  "signature":   (string)  Signature of comment (token+parentID+comment)
  "publickey":   (string)  Public key of user commenting
}

Response:
{
  "comment": {
    "token":        (string)  Censorship token
    "parentid":     (string)  Id of comment (defaults to '0' (top-level))
    "comment":      (string)  Comment
    "signature":    (string)  Signature of token+parentID+comment
    "publickey":    (string)  Public key of user 
    "commentid":    (string)  Id of the comment
    "receipt":      (string)  Server signature of the comment signature
    "timestamp":    (int64)   Received UNIX timestamp
    "censored":     (bool)    If comment has been censored
    "userid":       (string)  User id
    "username":     (string)  Username
  }
}`

type NewCommentCmd struct {
	Args struct {
		Token    string `positional-arg-name:"token" required:"true"`
		Comment  string `positional-arg-name:"comment" required:"true"`
		ParentID string `positional-arg-name:"parentID"`
	} `positional-args:"true"`
}

func (cmd *NewCommentCmd) Execute(args []string) error {
	token := cmd.Args.Token
	comment := cmd.Args.Comment
	parentID := cmd.Args.ParentID

	id := config.LoggedInUserIdentity
	if id == nil {
		return ErrNotLoggedIn
	}

	// Setup new comment request
	sig := id.SignMessage([]byte(token + parentID + comment))
	nc := &v1.NewComment{
		Token:     token,
		ParentID:  parentID,
		Comment:   comment,
		Signature: hex.EncodeToString(sig[:]),
		PublicKey: hex.EncodeToString(id.Public.Key[:]),
	}

	var ncr v1.NewCommentReply
	err := Ctx.Post(v1.RouteNewComment, nc, &ncr)
	if err != nil {
		return err
	}

	if !config.JSONOutput {
		if &ncr.Comment != nil {
			fmt.Printf("  %v\n", ncr.Comment.Token)
			fmt.Printf("  %v\n", ncr.Comment.Comment)
			fmt.Printf("      Submitted by: %v\n", ncr.Comment.Username)
			fmt.Printf("                at: %v\n",
				time.Unix(ncr.Comment.Timestamp, 0).String())
		}
	}
	return nil
}
