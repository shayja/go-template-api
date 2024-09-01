package utils

import "github.com/google/uuid"

func IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
 }
 
 func createNewUUID() uuid.UUID {
   id := uuid.Must(uuid.NewRandom())
   return id
 }
 