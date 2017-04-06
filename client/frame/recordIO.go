package frame

	
import  ( 
    "fmt"
    "time"
    "sync/atomic"
    "encoding/binary"
    "bufio"
    )

type RecordIO  struct {
        ctx context.Context
 }

type Bytes  []byte


//helper functionto read a line with bufio
func Readln(r  *bufio.Reader) (Bytes, error) {
    var (
        isPrefix bool = true
        err error = nil
        line, ln []byte
    )
    for isPrefix && err == nil {
        line, isPrefix, err = r.ReadLine()
        ln = append(ln, line...)
    }
    return Bytes,err
}

//helper function to check whether a flag is set or not in atomic sense
func isAtomicValSet( v *uint32 ) bool
{
    opsFinal := atomic.LoadUint32(v)
    if opsFinal > 0 {
        return true;
    }
    else {
        return false;
    }
}
    
//read buffer size (might be set to lower values )
const BUF_SIZE=1024

//interface function for Read
func (*RecordIO )  Read(reader io.ReadCloser, f FrameRead)  (interface{}, error) {

    //flag to be set when context is 
    var isContextDone  uint32 =0;
    
    errc:=make(chan error) //channel for sending error out
    
    //initialize the bufio
    r:=bufio.NewReader(reader)
    
    //to be able to wait for the graceful exit of below function
    var wg sync.WaitGroup
    wg.Add(1)
    //note: check how channels receive arrays, if they receive as copies, is it possible to send a channel an address of bytes from a local variable? if so, use the pointer version instead of copy by value
    
    // go function 
    //go func (r bufio.Reader,  rfc chan bytes[], erc chan error, ops * uint32,)   {
    go func ()   {
        //graceful exit
        defer wg.Done();
        
        //check if context is already cancelled/done, if so, return
        fcc:=func  () bool
        {
            opsFinal := atomic.LoadUint32(&isContextDone)
            if opsFinal > 0 {
                return true;
            }
            else {
                return false;
            }   
        }
        
        
     for ! fcc()  {
         b, err := Readln(r)
      
         
         //check if context is already cancelled/done, if so, return
         if fcc() {    
            return;
        }
         if err!= nil {
               errc <- err  
               return;
        }     
         
         buf := bytes.NewBuffer(b) // b is []byte
         s, err := binary.ReadVaruint(buf)
         if err!= nil {
            //note/add: convert the error that byte is not read
            errc <- err         
            return;
        }
        
         //add some check constraints here to the read buffer size!!!!!!!!!!!  might be illegitimate! 
         
         //check for panics of memory alloc later
         
         trb = make([] byte, s,s)
         trbr= trb[:]
         l:= 0 //length of read bytes
         for !fcc()  {
                    (n , err )= r.ReadBytes(trbr);
                     l+=n
                    //check if context is already cancelled/done, if so, return
                    if fcc() {    
                        return;
                    }
             
             //read entire data
             // add check if received n is smaller and or n is larger but eof returned
             if err==io.EOF {
                 if l<s {
                     
                     errc<-error.NewError("channel closed before receiving frame")
                     return;
                 } else {
                     trbTmp = make([]byte, s,s)
                     copy(trbTmp,trb)
                     f(&trbTmp)
                    
                     errc <- error.NewError("channel closed")
                    return;  
                 }
                   
             }
             //if all entire frame read            
             else if l == s {
                trbTmp = make([]byte, s,s)
                copy(trbTmp,trb)
                 f(&trbTmp)
                 break;
                    
             } 
             //if entire frame is not read yet, continue reading in next loop
             else {
                 trbr=trb[l:]
            }
        }
    }
    }()
    done := ctx.Done()
   
         
     select {
     case msg := <-done:
         
         opsFinal := atomic.LoadUint32(&ops)
         atomic.AddUint32(&ops, 1)
         //cancel above loop
         wg.Wait();
         return (nil, ctx.err())
         
         case    e := <-errc:
         
         wg.Wait();
         errc.close();
         
         return (nil, e)
     }
} 

func (RecordIO *)  Write(writer io.WriteCloser, f  frameWritten)  (interface{}, error) {
   
        
        
        
 }
