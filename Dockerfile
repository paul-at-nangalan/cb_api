FROM debian

WORKDIR /service/bin

ADD ../../bin/cb_api /service/bin
EXPOSE 80 ### don't think we need this ... 

ENV  STATFILE=/service/out/datafile.csv
ENV  STATLOGINTERVAL=10s
#### Normally should avoid putting api keys into a repo - but OK for a demo
ENV  APIKEY="eyJhbGciOiJSUzI1NiIsImtpZCI6IkhKcDkyNnF3ZXBjNnF3LU9rMk4zV05pXzBrRFd6cEdwTzAxNlRJUjdRWDAiLCJ0eXAiOiJKV1QifQ.eyJhY2Nlc3NfdGllciI6InRyYWRpbmciLCJleHAiOjE5MzU4OTA5NDMsImlhdCI6MTYyMDUzMDk0MywianRpIjoiODZmMGI4OTMtMmM1MS00OGY5LTgxMjYtNDU2MDkyODc5ZTI4Iiwic3ViIjoiZmI2YTEwMzQtMmRlNC00ODU0LWE1NzctZGU4ZWNlYmMzZmJlIiwidGVuYW50IjoiY2xvdWRiZXQiLCJ1dWlkIjoiZmI2YTEwMzQtMmRlNC00ODU0LWE1NzctZGU4ZWNlYmMzZmJlIn0.FMnwkgMeMPlhcyphgIgpvSlB9ca___rZ0Lyu98PIgSPSZQ88t-AMU-ijGxFX25noJqK5DxN-VBlRreP_Fu6PSE8mxukE0xT8pcK3D5ezJXUQeeKPk6kIz33YvCNN90U95R2Ijwe4QA3ikkb4TiUaYVPgMMnav9XveBVPtE5X5k4ExE7mlX0BYQLnMDfCeChXsKevJHkZO5v5dzzfZ9CbCfI8lyqtXxIOxKj3topOsBMqkly_PZah8NNTmMxMpsxr6tRVI1T3mv05g8o8cdxJHXy-3RZCePKaJPFxYsKaeqoKJcdoSxESu0ZEqVh_1NuiimBvXifdysVnWcWkNpdh7Q"
ENV  NUM_EVENTPROCESSORS=2
ENV  STATSINTERVAL=1s
ENV  CHECKINTERVAL=1s

CMD ["/service/bin/cb_api" ]

