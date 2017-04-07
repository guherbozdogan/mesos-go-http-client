package frame

	
import  ( 
    "sync/atomic"
    "sync"
    binary "encoding/binary"
    "bufio"
    "context"
    "io"
    "bytes"
    "errors"
    )

type RecordIO  struct {
        
 }

func NewRecordIO()  FrameIO {
    return &RecordIO{}
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
    return ln,err
}

//helper function to check whether a flag is set or not in atomic sense
func isAtomicValSet( v *uint32 ) bool {
    opsFinal := atomic.LoadUint32(v)
    if opsFinal > 0 {
        return true;
    }    else {
        return false;
    }
}
    

//interface function for Read
func (rc RecordIO )  Read(ctx context.Context, reader io.ReadCloser, f FrameReadFunc, erf ErrorFunc)   {

    
    //flag to be set when context is done
    var icd  uint32 =0;
    
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
        fcc:=func  () bool    {
            opsFinal := atomic.LoadUint32(&icd)
            if opsFinal > 0 {
                return true;
            }  else {
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
            s, err := binary.ReadVarint(buf)
            if err!= nil {
                //note/add: convert the error that byte is not read
                errc <- err         
                return;
            }
        
            //add some check constraints here to the read buffer size!!!!!!!!!!!  might be illegitimate! 
         
            //check for panics of memory alloc later
         
            trb := make([] byte, s,s)
            trbr:= trb[:]
            l:= int64(0) //length of read bytes
            for !fcc()  {
                    ni , err  := r.Read(trbr)
                    n:=int64(ni)
                     l+=n
                    //check if context is already cancelled/done, if so, return
                    if fcc() {    
                        return;
                    }
             
                //read entire data
                // add check if received n is smaller and or n is larger but eof returned
                if err==io.EOF {
                    if l<s {
                     
                        errc<-errors.New("channel closed before receiving frame")
                         return;
                    } else {
                        trbTmp := make([]byte, s,s)
                        copy(trbTmp,trb)
                        f(ctx,Frame(trbTmp),s)
                    
                        errc <- errors.New("channel closed")
                        return;  
                     }
                   
                }    else if l == s {  //if all entire frame read            
                    trbTmp := make([]byte, s,s)
                    copy(trbTmp,trb)
                    f(ctx,Frame(trbTmp),s)
                     break;
                    
                } else {                //if entire frame is not read yet, continue reading in next loop
                    trbr=trb[l:]
                }
            }
        }
    }()
    
    done := ctx.Done()
   
         
     select {
     case  <-done:
         
         atomic.AddUint32(&icd, 1)
         //wait graceful close
         wg.Wait();
         erf(ctx, ctx.Err())
         close(errc)
         return;
         
         case    e := <-errc:
         
         wg.Wait();
         erf(ctx,e);
         close(errc);
         
         return 
     }
} 

func (*RecordIO )  Write(c context.Context,   writer io.WriteCloser, f  FrameWritten, erf ErrorFunc)   {
   
    //not need to implement right now:)
    //to be added later 
      return 
        
 }
