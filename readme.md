# Thanos - Load Balancer 
Thanos is a high-performance load balancer implemented in Golang. It efficiently distributes network traffic across multiple backend servers, ensuring optimal resource utilization and improved system reliability.
## Features 
- Efficient traffic distribution using round-robin, least connection algoritm. 
- Concurrent request handling for improved performance 
- Configurable backend server management 
- Logging and monitoring capabilities  
 
## Installation 

 1. Clone the repository:`

```
git clone https://github.com/hritesh04/Thanos.git
```

2. Navigate to the project directory:

```
cd Thanos
```

3. Build the project:

```
make install
```

## Configuration

Edit the `example.yml` file to set up your loadbalancer configuration :

```yaml
type: round-robin
port: 3000
backends:
  - http://localhost:3001
  - http://localhost:3002
```

## Usage 

```
thanos start --config=example.yml
```

## Contrbuting

Contributions are welcome! Please feel free to submit a Pull Request.