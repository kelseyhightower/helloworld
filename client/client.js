var PROTO_PATH = __dirname + '/../helloworld/helloworld.proto';
var grpc = require('grpc');
var pb = grpc.load(PROTO_PATH);

function main() {
  var client = new pb.helloworld.Greeter('127.0.0.1:10000', grpc.credentials.createInsecure());
  var user = "Kelsey"
  client.sayHello({name: user}, function(err, response) {
    console.log(response.message);
  });
}

main();
