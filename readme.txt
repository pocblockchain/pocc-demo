1。 本demo演示了如何在pocc测试网上构建、签名、发送转账交易的过程，3.113.3.207:26657 是测试网地址。
2。 请参考tests中的TestBasicSend()中的流程构建转账交易
c.send -> 构建bank.MsgSend消息 ->调用c.SignAndBroadCastTx去签名和发送交易
                                        |
                                        |->  SignMsg  获取ChainID, accountNumber, Sequence, MinGasPrice, Gas构建交易并用私钥签名,
                                        |
                                        |-> BroadcastTxSync()，发送交易到mempool

3. 测试网浏览器地址: https://explorer-test.pocblockchain.io/

