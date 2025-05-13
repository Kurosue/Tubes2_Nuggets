# Tugas Besar 2 IF2211 - Little Alchemy 2 Recipe Explorer!!
 
<div align="center">
  <img src="https://github.com/user-attachments/assets/50b6c2d2-170c-476a-b7e2-563448c902c5" alt="App Preview" />
</div>

 <br>
 <div align="center">
   <h3 align="center">Tech Stacks and Languages</h3>
 
   <p align="center">
 
[![Next.js](https://img.shields.io/badge/Next.js-000000?style=for-the-badge&logo=next.js&logoColor=white)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![React](https://img.shields.io/badge/React-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://reactjs.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwind-css&logoColor=white)](https://tailwindcss.com/)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
 
   </p>
 </div>

 <p align="center">
    <br />
    <a href="https://github.com/Kurosue/Tubes2_Nuggets/releases/">Releases</a>
    ·
    <a href="https://github.com/Kurosue/Tubes2_Nuggets/docs">Project Report (Bahasa Indonesia)</a>
</p>

 <div align="justify">  </div>
<br />
<br />

### Overview
<br />
Nuggets is inspired by Little Alchemy 2, allowing users to discover recipes for creating various elements through combinations. The application leverages parallel BFS and DFS algorithms implemented in Go to efficiently find multiple recipes. The frontend provides an interactive tree visualization built with D3.js, showing the relationships between elements and their ingredients. Users can compare the performance of different algorithms in real-time while exploring element combinations.

ACCESS THE WEB ON [HERE](https://nuggets.bwks.link/)
 ---
 
 ## Installation & Setup
 
 ### Requirements
 > - Git
 > - Golang
 > - npm
 > - Docker

### Dependencies
 > - Gin

 <br/>

 ### Installing Dependencies and Requirement

<a id="dependencies"></a>
1. Install [Golang](https://golang.org/doc/install/source) and [npm](https://docs.npmjs.com/downloading-and-installing-nodejs-and-npm)
2. Install [Docker](https://docs.docker.com/get-docker/)
3. Install Gin
   ```bash
   go get -u github.com/gin-gonic/gin
   ```
4. Install [Node.js](https://nodejs.org/en/download/) and [npm](https://docs.npmjs.com/downloading-and-installing-nodejs-and-npm)

<br>  
<br/>  


 ---
 ## How to Run
 1. Clone the repository
    ```   bash
    git clone https://github.com/Kurosue/Tubes2_Nuggets.git
    ```
 2. Go to the project directory:
    ```bash
    cd Tubes2-Nuggets
    ```
 3. Start the application using Docker Compose:
    ```bash
    docker-compose up
    ```
 4. Access The Web on [http://localhost:8888]

## Development Setup
1. Backend
    ```   bash
    cd src/backend
    go run element-api/server.go
    ```
2. Frontend
    ```   bash
    cd src/frontend
    npm install
    npm run dev
    ```
> [!Note]
> Make sure that all of the dependencies are already installed
 ---
 <!-- CONTRIBUTOR -->
 <div align="center" id="contributor">
   <strong>
     <h3> Nuggets </h3>
     <table align="center">
       <tr align="center">
         <td>NIM</td>
         <td>Name</td>
         <td>GitHub</td>
       </tr>
       <tr align="center">
         <td>13523028</td>
         <td>Muhamamd Aditya Rahmadeni</td>
         <td><a href="https://github.com/Kurosue">@Kurosue</a></td>
       </tr>
       <tr align="center">
         <td>13523045</td>
         <td>Nadhif Radityo Nugroho</td>
         <td><a href="https://github.com/NadhifRadityo">@NadhifRadityo</td>
       </tr>
       <tr align="center">
         <td>13523052</td>
         <td>Adhimas Aryo Bimo</td>
         <td><a href="https://github.com/Ryonlunar">@Ryonlunar</a></td>
       </tr>
     </table>
   </strong>
 </div>
 <br/>
 <br/>
 <br/>
 <br/>
 
 <div align="center">
Nuggets • © 2025
 </div>
