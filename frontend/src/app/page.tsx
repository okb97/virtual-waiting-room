"use client"

import {useEffect, useState} from 'react';

export default function Home() {
  const [ticketId, setTicketId] = useState<string | null>(null)
  const [position, setPosition] = useState<number | null>(null)
  const [waitTime, setWaitTime] = useState<number | null>(null)

  useEffect(() => {
    const joinQueue = async () => {
      console.log("sending POST /api/queue")
      const res = await fetch("/api/queue",{
        method : "POST"
      });
      console.log("response status:", res.status);
      const data = await res.json();
      console.log("response data:", data);
      setTicketId(data.ticketId)
      setPosition(5)
      setWaitTime(100)
    }
    joinQueue();
  },[])

  useEffect(() => {
    if(!ticketId) return

    const interval = setInterval(async() => {
      const res = await fetch(`/api/queue?ticketId=${ticketId}`)
      const data = await res.json()
      setPosition(data.position)
      setWaitTime(data.waitTime)
      if(data.position == 0){
        console.log("あなたの番が来ました！チケット購入ページへ移動します。")
        clearInterval(interval)

        const checkInRes = await fetch(`/api/checkin`,{
          method:'DELETE'
        })

        if(checkInRes.ok){
          window.location.href = '/purchase';
        }
        else{
          console.error("キューの削除に失敗しました")
        }
      }
    },30000)
    return () => clearInterval(interval)
  },[ticketId])


  return (
    <main>
      <h1>仮想待合室</h1>
      <p className="mb-2">あなたのチケットID: <strong>{ticketId}</strong></p>
          <p className="mb-2">現在の順番: {position} 番目</p>
          <p className="mb-4">推定待ち時間: {waitTime} 分</p>
    </main>
  );
}
