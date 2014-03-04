package query

import (
	"pilosa/db"
	"strings"

	"tux21b.org/v1/gocql/uuid"
)

type QueryInput interface{}

type QueryResults struct {
	Data interface{}
}

type PqlList []PqlListItem

type PqlListItem struct {
	Id    *uuid.UUID
	Label string
	PQL   string
}

type Query struct {
	Id         *uuid.UUID
	Operation  string
	Args       map[string]interface{}
	Subqueries []Query

	// deprecated:
	Inputs []QueryInput //"strconv"
	// Represents a parsed query. Inputs can be Query or Bitmap objects
	// Maybe Bitmap and Query objects should have different fields to avoid using interface{}
	ProfileId uint64 // used only for set() queries
	N         int    // TODO: I think we should make this a generic map for any attributes related to the query
}

func QueryPlanForPQL(database *db.Database, pql string, destination *db.Location) *QueryPlan {
	tokens, err := Lex(pql)
	if err != nil {
		panic(err)
	}
	return QueryPlanForTokens(database, tokens, destination)
}

func QueryForPQL(pql string) *Query {
	tokens, err := Lex(pql)
	if err != nil {
		panic(err)
	}
	return QueryForTokens(tokens)
}

func QueryForTokens(tokens []Token) *Query {
	query, err := Parse(tokens)
	if err != nil {
		panic(err)
	}
	return query
}

func QueryPlanForTokens(database *db.Database, tokens []Token, destination *db.Location) *QueryPlan {
	query := QueryForTokens(tokens)
	return QueryPlanForQuery(database, query, destination)
}

func QueryPlanForQuery(database *db.Database, query *Query, destination *db.Location) *QueryPlan {
	query_planner := QueryPlanner{Database: database}
	id := uuid.RandomUUID()
	query_plan := query_planner.Plan(query, &id, destination)
	return query_plan
}

func TokensToString(tokens []Token) string {
	var str []string
	for i, _ := range tokens {
		str = append(str, tokens[i].Text)
	}
	// for now, we're just using this function to pull the filter out of the outer function "outerfunc(filter)"
	str = str[2 : len(str)-1]
	return strings.Join(str, "")
}
