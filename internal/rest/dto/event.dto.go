package dto

type EventCreateReq struct {
	Name string `json:"name" form:"name"`
}

type EventEditReq struct {
	Name              string `json:"name" form:"name"`
	Desc              string `json:"desc" form:"desc"`
	GroupMemberNum    int    `json:"group_member_num" form:"group_member_num"`
	ParticipantTarget int    `json:"participant_target" form:"participant_target"`
	Period            string `json:"period" form:"period"`
}
