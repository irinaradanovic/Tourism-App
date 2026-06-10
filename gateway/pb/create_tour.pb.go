// Minimal protobuf definitions for gateway -> tours TourCommandService calls.
package pb

import protoimpl "google.golang.org/protobuf/runtime/protoimpl"

type Difficulty int32

const (
	Difficulty_DIFFICULTY_UNSPECIFIED Difficulty = 0
	Difficulty_EASY                   Difficulty = 1
	Difficulty_MEDIUM                 Difficulty = 2
	Difficulty_HARD                   Difficulty = 3
)

type TransportType int32

const (
	TransportType_TRANSPORT_UNSPECIFIED TransportType = 0
	TransportType_WALKING               TransportType = 1
	TransportType_BICYCLE               TransportType = 2
	TransportType_CAR                   TransportType = 3
)

type CreateTourRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Difficulty    Difficulty             `protobuf:"varint,3,opt,name=difficulty,proto3,enum=tours.Difficulty" json:"difficulty,omitempty"`
	Tags          []string               `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	AuthorId      int64                  `protobuf:"varint,5,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Durations     []*TourDurationProto    `protobuf:"bytes,6,rep,name=durations,proto3" json:"durations,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTourRequest) Reset()         { *x = CreateTourRequest{} }
func (x *CreateTourRequest) String() string { return x.Title }
func (x *CreateTourRequest) ProtoMessage()  {}

type PublishTourRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TourId        string                 `protobuf:"bytes,1,opt,name=tour_id,json=tourId,proto3" json:"tour_id,omitempty"`
	AuthorId      int64                  `protobuf:"varint,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Role          string                 `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PublishTourRequest) Reset()         { *x = PublishTourRequest{} }
func (x *PublishTourRequest) String() string { return x.TourId }
func (x *PublishTourRequest) ProtoMessage()  {}

type TourDurationProto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TransportType TransportType          `protobuf:"varint,1,opt,name=transport_type,json=transportType,proto3,enum=tours.TransportType" json:"transport_type,omitempty"`
	Minutes       int32                  `protobuf:"varint,2,opt,name=minutes,proto3" json:"minutes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TourDurationProto) Reset()         { *x = TourDurationProto{} }
func (x *TourDurationProto) String() string { return "" }
func (x *TourDurationProto) ProtoMessage()  {}

type TourProto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AuthorId      int64                  `protobuf:"varint,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Title         string                 `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Difficulty    Difficulty             `protobuf:"varint,5,opt,name=difficulty,proto3,enum=tours.Difficulty" json:"difficulty,omitempty"`
	Tags          []string               `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	Price         float64                `protobuf:"fixed64,7,opt,name=price,proto3" json:"price,omitempty"`
	Status        string                 `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	DistanceKm    float64                `protobuf:"fixed64,9,opt,name=distance_km,json=distanceKm,proto3" json:"distance_km,omitempty"`
	PublishedAt   string                 `protobuf:"bytes,10,opt,name=published_at,json=publishedAt,proto3" json:"published_at,omitempty"`
	ArchivedAt    string                 `protobuf:"bytes,11,opt,name=archived_at,json=archivedAt,proto3" json:"archived_at,omitempty"`
	Durations     []*TourDurationProto    `protobuf:"bytes,12,rep,name=durations,proto3" json:"durations,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TourProto) Reset()         { *x = TourProto{} }
func (x *TourProto) String() string { return x.Title }
func (x *TourProto) ProtoMessage()  {}

type CreateTourResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Tour          *TourProto             `protobuf:"bytes,1,opt,name=tour,proto3" json:"tour,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTourResponse) Reset()         { *x = CreateTourResponse{} }
func (x *CreateTourResponse) String() string { return "" }
func (x *CreateTourResponse) ProtoMessage()  {}

type TourCommandResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Tour          *TourProto             `protobuf:"bytes,1,opt,name=tour,proto3" json:"tour,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TourCommandResponse) Reset()         { *x = TourCommandResponse{} }
func (x *TourCommandResponse) String() string { return x.Message }
func (x *TourCommandResponse) ProtoMessage()  {}
