package db

func (db *DatabaseConnection) AddServer(gID string) error {
	return db.Create(&PugServer{ServerID: gID}).Error
}

func (db *DatabaseConnection) RemoveServer(gID string) error {
	return db.Where("server_id = ?", gID).Delete(&PugServer{}).Error
}

func (db *DatabaseConnection) AddChannel(cID, gID string) error {
	var parentServer PugServer
	db.Where("server_id = ?", gID).First(&parentServer)
	return db.Create(&QueueChannel{ChanID: cID, PugServerID: parentServer.ID}).Error
}

func (db *DatabaseConnection) RemoveChannel(cID string) error {
	return db.Where("chan_id = ?", cID).Delete(&QueueChannel{}).Error
}

func (db *DatabaseConnection) GetRegisteredIds() (servers []*PugServer, queues []*QueueChannel) {
	db.Find(&servers)
	db.Find(&queues)
	return
}
