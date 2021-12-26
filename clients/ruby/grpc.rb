require './table_pb'
require 'grpc'
require './table_services_pb'

if ARGV.length != 4
    puts 'wrong number of arguments supplied, should be 4, serverIP, tableID1, tableID2, dbID'
    return
end


host = ARGV[0]
id1 = ARGV[1]
id2 = ARGV[2]
dbId = ARGV[3]

stub = TableProductProvider::Stub.new(host + ':9090', :this_channel_is_insecure)

# stub.GetTableProduct(Filter.new('6162cf849d5b528868c1b14d', '6162c07fd22ed745bb63611c', '6162bfa4d22ed745bb63611b'))
res = stub.get_table_product(Filter.new(id1: id1,id2: id2, dbId:dbId))

puts res