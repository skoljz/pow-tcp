### **TCP Server with Proof-of-Work DDoS protect (simplified hashcash)**

```/server``` - **TCP Server** 

```/client``` - **CLI app for server**

____

### **Usage**

#### 1. Run
```bash
docker-compose up -d # up containers
docker-compose exec pow-client ./app # enter to cli
```

#### 2. Commands for CLI
```req``` - send request to TCP server and solve work
```q, quit``` - quit from CLI


#### 3. Redis usage (optional)

For redis usage you need add quotes with **RPUSH** / **LPUSH** with '**quotes**' key before request
______

**For test run use ```make test``` in ```/server``` or ```/cmd``` directory**
