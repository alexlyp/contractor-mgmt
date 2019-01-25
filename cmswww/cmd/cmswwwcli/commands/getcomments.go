package commands

import (
	"fmt"
	"time"

	"github.com/decred/contractor-mgmt/cmswww/api/v1"
	"github.com/decred/contractor-mgmt/cmswww/cmd/cmswwwcli/config"
)

// Help message displayed for the command 'politeiawwwcli help getcomments'
var GetCommentsCmdHelpMsg = `getcomments "token" 

Get comments for a proposal.

Arguments:
1. token       (string, required)   Proposal censorship token

Result:
{
  "comments": [
    {
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
  ]
}`

type GetCommentsCmd struct {
	Args struct {
		Token string `positional-arg-name:"token"`
	} `positional-args:"true" required:"true"`
}

func (cmd *GetCommentsCmd) Execute(args []string) error {
	gc := v1.GetComments{
		Token: cmd.Args.Token,
	}
	var gcr v1.GetCommentsReply
	err := Ctx.Get(v1.RouteComments, gc, &gcr)
	if err != nil {
		return err
	}

	if !config.JSONOutput {
		for _, c := range gcr.Comments {
			fmt.Printf("  %v\n", c.Token)
			fmt.Printf("  %v\n", c.Comment)
			fmt.Printf("      Submitted by: %v\n", c.Username)
			fmt.Printf("                at: %v\n",
				time.Unix(c.Timestamp, 0).String())
		}
	}

	return nil
}
