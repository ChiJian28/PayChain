import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Alert, Button, Card, Col, Divider, Flex, Form, Input, InputNumber, Layout, List, Row, Space, Statistic, Tag, message } from 'antd'
import { Block, Transaction, getBalance, getBlockchain, getPending, postTransfer, postFaucet } from '../lib/api'
import { useState } from 'react'
import { useAppStore } from '../store/useAppStore'

const { Header, Content } = Layout

function BlockCard({ block }: { block: Block }) {
  return (
    <Card
      size="small"
      title={`#${block.Index}`}
      extra={<span className="text-xs">{new Date(block.Timestamp * 1000).toLocaleString()}</span>}
      className="mb-2"
    >
      <Space direction="vertical" size={4} className="w-full">
        <div className="text-xs break-all">Hash: <code>{block.Hash.slice(0, 16)}…</code></div>
        <div className="text-xs break-all">Prev: <code>{block.PrevHash.slice(0, 16)}…</code></div>
        <div className="text-xs">Tx: {block.Transactions?.length ?? 0}</div>
      </Space>
    </Card>
  )
}

export default function App() {
  const qc = useQueryClient()
  const [user, setUser] = useState('alice')
  const { from, to, amount, setFrom, setTo, setAmount } = useAppStore()

  const blockchain = useQuery({ queryKey: ['chain'], queryFn: getBlockchain, refetchInterval: 3000 })
  const pending = useQuery({ queryKey: ['pending'], queryFn: getPending, refetchInterval: 1000 })
  const balance = useQuery({ queryKey: ['balance', user], queryFn: () => getBalance(user), refetchInterval: 3000 })

  const transfer = useMutation({
    mutationFn: postTransfer,
    onSuccess: () => {
      message.success('Queued successfully')
      qc.invalidateQueries({ queryKey: ['pending'] })
    },
    onError: (e: any) => message.error(e?.message || 'Submission failed')
  })

  const faucet = useMutation({
    mutationFn: () => postFaucet({ to: 'alice', amount: 1000 }),
    onSuccess: (data) => {
      message.success(`Minted to ${data.user}, balance: ${data.balance}`)
      qc.invalidateQueries({ queryKey: ['balance', 'alice'] })
    },
    onError: (e: any) => message.error(e?.message || 'Faucet failed')
  })

  return (
    <Layout className="min-h-screen">
      <Header className="bg-white border-b">
        <div className="max-w-6xl mx-auto text-lg font-semibold">PayChain Dashboard</div>
      </Header>
      <Content>
        <div className="max-w-6xl mx-auto p-4">
          <Row gutter={[16, 16]}>
            <Col xs={24} md={12}>
              <Card title="Transfer">
                <Form layout="vertical" onFinish={() => transfer.mutate({ from, to, amount })}>
                  <Form.Item label="From">
                    <Input value={from} onChange={e => setFrom(e.target.value)} />
                  </Form.Item>
                  <Form.Item label="To">
                    <Input value={to} onChange={e => setTo(e.target.value)} />
                  </Form.Item>
                  <Form.Item label="Amount">
                    <InputNumber min={1} value={amount} onChange={v => setAmount(Number(v || 0))} />
                  </Form.Item>
                  <Button type="primary" htmlType="submit" loading={transfer.isPending}>Submit</Button>
                </Form>
              </Card>

              <Divider />

              <Card title="Check Balance">
                <Flex gap={8} align="center">
                  <Input style={{ width: 180 }} value={user} onChange={e => setUser(e.target.value)} />
                  <Button onClick={() => qc.invalidateQueries({ queryKey: ['balance', user] })}>Refresh</Button>
                  <Button type="dashed" onClick={() => faucet.mutate()} loading={faucet.isPending}>Give Alice +1000</Button>
                </Flex>
                <div className="mt-3">
                  <Statistic title={balance.data?.user || user} value={balance.data?.balance ?? 0} />
                </div>
              </Card>
            </Col>

            <Col xs={24} md={12}>
              <Card title={<Flex align="center" gap={8}>Blockchain<Tag color="blue">{blockchain.data?.length ?? 0}</Tag></Flex>}>
                <div className="max-h-96 overflow-auto">
                  <List
                    dataSource={[...(blockchain.data || [])].reverse()}
                    renderItem={(b) => <BlockCard key={b.Hash} block={b} />}
                  />
                </div>
              </Card>

              <Divider />

              <Card title={<Flex align="center" gap={8}>Pending<Tag color="orange">{pending.data?.length ?? 0}</Tag></Flex>}>
                {pending.data && pending.data.length === 0 && (
                  <Alert type="info" message="No pending transactions" showIcon />
                )}
                <List
                  dataSource={pending.data || []}
                  renderItem={(tx: Transaction, idx) => (
                    <List.Item>
                      <Flex wrap gap={8} className="text-xs">
                        <Tag color="default">#{idx + 1}</Tag>
                        <span>From: {tx.From || '-'}</span>
                        <span>To: {tx.To || '-'}</span>
                        <span>Amount: {tx.Amount}</span>
                        <span>Time: {new Date(tx.Time * 1000).toLocaleString()}</span>
                      </Flex>
                    </List.Item>
                  )}
                />
              </Card>
            </Col>
          </Row>
        </div>
      </Content>
    </Layout>
  )
}








