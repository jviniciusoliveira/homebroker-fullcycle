export type Asset = {
  id: string
  symbol: string
  price: number
}

export type WalletAsset = {
  id: string
  wallet_id: string
  asset_id: string
  shares: number
  Asset: Asset
}

export type Order = {
  id: string
  wallet_id: string
  asset_id: string
  shares: string
  partial: string
  price: string
  type: 'BUY' | 'SELL'
  created_at: string
  updated_at: string
  status: 'PENDING' | 'OPEN' | 'CLOSED' | 'FAILED'
  Asset: Pick<Asset, 'id' | 'symbol'>
}
