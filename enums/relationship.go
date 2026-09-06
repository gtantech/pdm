package enums

type RelationshipType string

const (
	SS RelationshipType = "SS" // start to start relationship
	FS RelationshipType = "FS" // finish to start relationship
	SF RelationshipType = "SF" // start to finish relationship
	FF RelationshipType = "FF" //finish to finish relationship
)
