package server


type topic struct{
    name string
}

type Broker struct{
    ip string
    port int32 
    topics []topic
}

func NewBroker(){
    return &Broker{
        
    }
}