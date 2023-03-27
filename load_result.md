> k6 run load.js


          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: tests/load.js
     output: -

  scenarios: (100.00%) 2 scenarios, 130 max VUs, 1m0s max duration (incl. graceful stop):
           * consume: 60 looping VUs for 30s (exec: consume, gracefulStop: 30s)
           * queueAdd: 70 looping VUs for 30s (exec: queueAdd, gracefulStop: 30s)


running (0m30.0s), 000/130 VUs, 2132046 complete and 0 interrupted iterations
consume  ✓ [======================================] 60 VUs  30s
queueAdd ✓ [======================================] 70 VUs  30s

     ✓ status was 200

     checks.........................: 100.00% ✓ 2132046      ✗ 0      
     data_received..................: 205 MB  6.8 MB/s
     data_sent......................: 222 MB  7.4 MB/s
     http_req_blocked...............: avg=1.77µs  min=0s      med=0s     max=41.93ms  p(90)=1µs    p(95)=1µs   
     http_req_connecting............: avg=190ns   min=0s      med=0s     max=26.46ms  p(90)=0s     p(95)=0s    
     http_req_duration..............: avg=1.28ms  min=17µs    med=791µs  max=102.56ms p(90)=2.59ms p(95)=3.89ms
       { expected_response:true }...: avg=1.28ms  min=17µs    med=791µs  max=102.56ms p(90)=2.59ms p(95)=3.89ms
     http_req_failed................: 0.00%   ✓ 0            ✗ 2132046
     http_req_receiving.............: avg=15.25µs min=1µs     med=4µs    max=57.27ms  p(90)=6µs    p(95)=11µs  
     http_req_sending...............: avg=6.74µs  min=1µs     med=2µs    max=70.19ms  p(90)=4µs    p(95)=6µs   
     http_req_tls_handshaking.......: avg=0s      min=0s      med=0s     max=0s       p(90)=0s     p(95)=0s    
     http_req_waiting...............: avg=1.26ms  min=13µs    med=780µs  max=102.55ms p(90)=2.56ms p(95)=3.83ms
     http_reqs......................: 2132046 71064.386211/s
     iteration_duration.............: avg=1.71ms  min=30.37µs med=1.06ms max=102.58ms p(90)=3.52ms p(95)=5.3ms 
     iterations.....................: 2132046 71064.386211/s
     vus............................: 130     min=130        max=130  
     vus_max........................: 130     min=130        max=130  
