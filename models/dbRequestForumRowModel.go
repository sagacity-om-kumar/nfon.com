package models

type DBResponseDoc struct {
	MsgID   int    `db:"TICKET_MESSAGE_UID" json:"msgid"`
	DocID   int    `db:"TICKET_MSG_DOCUMENT_UID" json:"docid"`
	Doc     []byte `db:"TICKET_MSG_DOC_FILE_STORAGE" json:"doc"`
	DocName string `db:"TICKET_MSG_DOC_FILE_NAME" json:"docname"`
}
