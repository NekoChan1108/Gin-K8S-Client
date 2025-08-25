# WebSocket 模块设计说明：
## 🔄 整体流程 前端输入 → 后端处理
```mermaid
graph LR
    A[前端] -->|WebSocket发送数据| B[ReadLoop]
    B -->|读取数据| C[inChan]
    C -->|获取数据| D[ReadMsg]
    D -->|业务逻辑处理| E[应用层]

```
1. 前端通过 WebSocket 发送数据
2. ReadLoop 读取数据 → 放入 inChan
3. ReadMsg() 获取数据供业务逻辑处理
## 🔄 整体流程 后端响应 → 前端展示

```mermaid
graph LR
    A[应用层] -->|需要返回数据| B[WriteMsg]
    B -->|放入数据| C[outChan]
    C -->|取出数据| D[WriteLoop]
    D -->|WebSocket发送| E[前端]

```
1. 后端业务逻辑需要返回数据
2. WriteMsg() 将数据放入 outChan
3. WriteLoop 取出数据 → 通过 WebSocket 发送给前端

## 🏗️ 架构设计
```mermaid
graph TD
    subgraph WebSocket连接层
        A[WebSocket连接]
    end
    
    subgraph 读循环
        B[ReadLoop]
        C[inChan]
        D[ReadMsg]
    end
    
    subgraph 写循环
        E[WriteLoop]
        F[outChan]
        G[WriteMsg]
    end
    
    subgraph 控制层
        H[closeChan]
        I[Close]
    end
    
    A -->|读取| B
    B --> C
    C --> D
    
    G --> F
    F --> E
    E -->|写入| A
    
    H --> B
    H --> D
    H --> E
    H --> G
    I --> H

```