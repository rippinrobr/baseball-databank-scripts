package databank


import (
  "os"
  "log"
  "errors"
  "path/filepath"

  "github.com/rippinrobr/baseball-stats-db/internal/platform/parsers/csv"

  "github.com/rippinrobr/baseball-stats-db/internal/platform/db"

)

// BattingPost is a model that maps the CSV to a DB Table
type BattingPost struct {
   Yearid   int `json:"yearID"  csv:"yearID"  db:"yearID"  bson:"yearID"`
   Round   string `json:"round"  csv:"round"  db:"round,omitempty"  bson:"round"`
   Playerid   string `json:"playerID"  csv:"playerID"  db:"playerID,omitempty"  bson:"playerID"`
   Teamid   string `json:"teamID"  csv:"teamID"  db:"teamID,omitempty"  bson:"teamID"`
   Lgid   string `json:"lgID"  csv:"lgID"  db:"lgID,omitempty"  bson:"lgID"`
   G   int `json:"g"  csv:"G"  db:"G"  bson:"G"`
   Ab   int `json:"aB"  csv:"AB"  db:"AB"  bson:"AB"`
   R   int `json:"r"  csv:"R"  db:"R"  bson:"R"`
   H   int `json:"h"  csv:"H"  db:"H"  bson:"H"`
   Doubles   int `json:"doubles"  csv:"2B"  db:"doubles"  bson:"doubles"`
   Triples   int `json:"triples"  csv:"3B"  db:"triples"  bson:"triples"`
   Hr   int `json:"hR"  csv:"HR"  db:"HR"  bson:"HR"`
   Rbi   int `json:"rBI"  csv:"RBI"  db:"RBI"  bson:"RBI"`
   Sb   int `json:"sB"  csv:"SB"  db:"SB"  bson:"SB"`
   Cs   int `json:"cS"  csv:"CS"  db:"CS"  bson:"CS"`
   Bb   int `json:"bB"  csv:"BB"  db:"BB"  bson:"BB"`
   So   int `json:"sO"  csv:"SO"  db:"SO"  bson:"SO"`
   Ibb   int `json:"iBB"  csv:"IBB"  db:"IBB"  bson:"IBB"`
   Hbp   int `json:"hBP"  csv:"HBP"  db:"HBP"  bson:"HBP"`
   Sh   int `json:"sH"  csv:"SH"  db:"SH"  bson:"SH"`
   Sf   int `json:"sF"  csv:"SF"  db:"SF"  bson:"SF"`
   Gidp   int `json:"gIDP"  csv:"GIDP"  db:"GIDP"  bson:"GIDP"`
  inputDir string
}

// GetTableName returns the name of the table that the data will be stored in
func (m *BattingPost) GetTableName() string {
  return "battingpost"
}

// GetFileName returns the name of the source file the model was created from
func (m *BattingPost) GetFileName() string {
  return "BattingPost.csv"
}

// GetFilePath returns the path of the source file
func (m *BattingPost) GetFilePath() string {
  return filepath.Join(m.inputDir, "BattingPost.csv")
}

// SetInputDirectory sets the input directory's path so it can be used to create the full path to the file
func (m *BattingPost) SetInputDirectory(inputDir string) {
  m.inputDir=inputDir
}

// GenParseAndStoreCSV returns a function that will parse the source file,\n//create a slice with an object per line and store the data in the db
func (m *BattingPost) GenParseAndStoreCSV(f *os.File, repo db.Repository, pfunc csv.ParserFunc) (ParseAndStoreCSVFunc, error) {
  if f == nil {
    return func() error{return nil}, errors.New("nil File")
  }

  return func() error {
    rows := make([]*BattingPost, 0)
    numErrors := 0
    err := pfunc(f, &rows)
    if err == nil {
       numrows := len(rows)
       if numrows > 0 {
         log.Println("BattingPost ==> Truncating")
         terr := repo.Truncate(m.GetTableName())
         if terr != nil {
            log.Println("truncate err:", terr.Error())
         }

         log.Printf("BattingPost ==> Inserting %d Records\n", numrows)
         for _, r := range rows {
           ierr := repo.Insert(m.GetTableName(), r)
           if ierr != nil {
             log.Printf("Insert error: %s\n", ierr.Error())
             numErrors++
           }
         }
       }
       log.Printf("BattingPost ==> %d records created\n", numrows-numErrors)
    }

    return err
   }, nil
}
