package frame

	
import "fmt"
import "time"
import "sync/atomic"
import "encoding/binary"

type RecordIO  struct {
        ctx context.Context
 }

type Bytes []byte

func Readln(r *bufio.Reader) (Bytes, error) {
        var (isPrefix bool = true
        err error = nil
            line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return Bytes,err
}

func isAtomicValSet( v uint64 *) bool
{
       opsFinal := atomic.LoadUint64(v)
        if opsFinal > 0 {
                return true;
        }
      else {
        return false;
    }
}
        
const BUF_SIZE=1024
        

func (RecordIO *)  Read(reader io.ReadCloser, f FrameRead)  (interface{}, error) {

    var isContextDone  uint64 =0;
    
    rfc := make(chan bytes[])  //channel for sending bytes out
    errC:=make(chan error) //channel for sending error out
    
    r:=bufio.NewReader(reader)
     
    //check how channels receive arrays, if they receive as copies, is it possible to send a channel an address of bytes from a local variable? if so, use the pointer version instead of copy by value
    go func (r bufio.Reader,  rfc chan bytes[], ops * uint64 , erc chan error)   {
        
        //check if context is already cancelled/done, if so, return
        f:=func  (rfc chan bytes[], ops * uint64 , erc chan error) bool
        {
            if isAtomicValSet(ops) {    
            rfc,close();
            erc.close();
            return true;
            }
            else  return false;
        }
        
        
     for ! f(rfc, ops, erc)  {
         
         b, err= Readln(r)
      
         //check if context is already cancelled/done, if so, return
         if f(rfc, ops, erc) {    
            return;
        }
              
         
         buf := bytes.NewBuffer(b) // b is []byte
         s, err := binary.ReadVaruint(buf)
         if err!= nil {
               erc <- err         
        }
        //add some check constraints here to the read buffer size!!!!!!!!!!!  might be illegitimate! 
         
         //check for panics of memory alloc later
         trb = make([] byte, BUF_SIZE,BUF_SIZE)
         trbr= trb[:]
         for {
             
                 //check if context is already cancelled/done, if so, return
                    if f(rfc, ops, erc) {    
                        return;
                    }
             
                    (n , err )= r.ReadBytes(trbr);
                    //check if context is already cancelled/done, if so, return
                    if f(rfc, ops, erc) {    
                        return;
                    }
             
             if err==io.EOF {
                 
                   rfc <- trb
                 return;
             }
             
             else if n == BUF_SIZE {
                   //check for panics of memory alloc later, dont know how to do now
                   trbTmp = make([]byte, len(trb), cap(trb) + BUF_SIZE)
                     copy(trbTmp,trb)
                    trb=trbTmp;
                   trbr=trb[len(trb):]
             } else {
                 trbr=trb[len(trb):]
                 
             }
         }
     }
    }(reader, f, c, rfc, isDone)
    
        
     done := ctx.Done()
   
         
     select {
     case msg := <-done:
         
         opsFinal := atomic.LoadUint64(&ops)
         
          atomic.AddUint64(&ops, 1)
         //cancel above loop
         
         case  
     }
    }
	 case msg := <- result
          //go call again the reader
         
        

    
        
        
        
    } (reader, f)
    
    
    
    
}

func (RecordIO *)  Write(writer io.WriteCloser, f  frameWritten)  (interface{}, error) {
   
        
        
        
 }
