# clone and build generation tool. Output omitted.
git clone git@github.com:gunnarmorling/1brc.git
./mvnw clean verify
./create_measurements.sh 1000000000
