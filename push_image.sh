TAG=a0.1.6

docker push gniang/drone-teams
docker tag gniang/drone-teams:latest gniang/drone-teams:$TAG
docker push gniang/drone-teams:$TAG
