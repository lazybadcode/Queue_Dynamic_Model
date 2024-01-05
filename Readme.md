NOTE:
- priority ของ queue per day เรียงจาก วันหยุดเสาร์อาทิตย์ วันหยุดพิเศษ วันธรรมดา
- checking config/config.yaml before use
- unittest using `mockery`, when add new method in usecase, you should run `make genmock-usecase`
  - install mockery on mac `brew install mockery`
- run unit test with `go test ./...` or `make test`

Feature:
- graphql
  - create queue
  - query queue (all, by queue id, by user id)
  - update queue
  - delete queue
- batch
  - send sms at 22:00 (only log in console)

HOW TO USE
1. start database
```shell
 docker compose up
```
2. start api (ENV = local,sit,uat,prod)
```shell
TZ=Asia/Bangkok ENV=local go run .
```
3. add index to mongo for avoid race condition
```javascript
db.queues.createIndex({ date: 1, no: 1 }, { unique: true })
```

4. Go to http://localhost:8080/ for test

Create
```graphql
mutation {
  createQueue(
    idCard: "12345678910"
    mobileNo: "0801234567"
    input: {
      note: "This is my note"
      testBoolean:true
      testFloat:8.88
      user: { name: "this is my name" }
    }
  ) {
    _id
    slot
  }
}

```

Query All
```graphql
query {
  queue(
    input: {
      _id: ""
      userId:""
    }
  ) {
    _id
    userId
    date
    testFloat
    testBoolean
    user {
      _id
      name
      idCard
      createdTime
    }
  }
}

```
Query 1 User
```graphql
query {
  queue(
    input: {
      _id: ""
      userId:"658aed573682a9a63ddb8a6e"
    }
  ) {
    _id
    userId
    date
    user {
      _id
      name
      idCard
      createdTime
    }
  }
}
```
Query By Date
```graphql
query {
  queue(
    input: {
      # _id: "658aed663682a9a63ddb8a6f"
      # userId:"658aed573682a9a63ddb8a6e"
      date:"20231226"
    }
  ) {
    _id
    userId
    date
    user {
      _id
      name
      idCard
      createdTime
    }
  }
}

```
Update Queue (specific date,slot,note) 
```graphql
mutation {
  updateQueue(
    id: "657d6f1676863b9c94c22242"
    date: "20231226"
    slot: 3
    input: {
      note: "This is my note"
      user: { name: "this is my name 2"}
    }
  ) {
    _id
    no
    slot
    note
    date
    user {
      name
      _id
    }
  }
}

```

Update Queue only note
```graphql
mutation {
  updateQueue(
    id: "657d6f1676863b9c94c22242"
    date: ""
    slot: 0
    input: {
      note: "This is my note"
      user: { name: "this is my name 2" }
    }
  ) {
    _id
    no
    slot
    note
    date
    user {
      name
      _id
    }
  }
}
```

Delete queue
```graphql
mutation {
  deleteQueue(
    id: "658aea00f9b1f054e7333519"
  ) {
    _id
    no
    slot
    note
    date
    user {
      name
      _id
    }
  }
}

```



mock special days
```javascript
db = db.getSiblingDB('test');

var specialDays = [
    { "date": "20230101", "description": "New Year's Day" },
    { "date": "20230208", "description": "Makha Bucha" },
    { "date": "20230406", "description": "Chakri Day" },
    { "date": "20230413", "description": "Songkran Festival" },
    { "date": "20230414", "description": "Songkran Festival" },
    { "date": "20230415", "description": "Songkran Festival" },
    { "date": "20230501", "description": "Labor Day" },
    { "date": "20230504", "description": "Coronation Day" },
    { "date": "20230506", "description": "Vesak Day" },
    { "date": "20230603", "description": "H.M. Queen Suthida's Birthday" },
    { "date": "20230703", "description": "Asalha Bucha Day" },
    { "date": "20230704", "description": "Buddhist Lent Day" },
    { "date": "20230728", "description": "H.M. King's Birthday" },
    { "date": "20230812", "description": "Mother's Day" },
    { "date": "20231023", "description": "Chulalongkorn Day" },
    { "date": "20231205", "description": "Father's Day" },
    { "date": "20231210", "description": "Constitution Day" },
    { "date": "20231231", "description": "New Year's Eve" }
];

db.special_days.insertMany(specialDays);
```

# ขั้นตอนในการเพิ่ม Query และ Mutation
หากคุณต้องการเพิ่ม Query หรือ Mutation ในโครงการของคุณ กรุณาปฏิบัติตามขั้นตอนดังนี้:
1. ไปที่ไฟล์ `schema/schema.go`: เพิ่มฟังก์ชันใน `type Query` หรือ `type Mutation` ตามที่ต้องการ อ้างอิงจากของเดิมได้เลยครับ
2. คุณจะต้อง implement logic สำหรับแต่ละฟังก์ชันที่คุณเพิ่มเข้ามาใหม่ โดยทั่วไปแล้วเรียก function ที่ implement logic ที่ package ชื่อ usecase
