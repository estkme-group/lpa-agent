export interface Information {
  readonly eid: string
  readonly default_smds: string
  readonly default_smdp?: string
}

export interface Profile {
  readonly iccid: string
  readonly isdpaid: string
  readonly profileName: string
  profileNickname?: string
  profileState: number
  readonly profileClass: number
  readonly serviceProviderName: string
}
