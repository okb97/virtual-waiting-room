"use client"

import {useEffect, useState} from 'react';

export default function Home() {
  const [ticketId, setTicketId] = useState<string | null>(null)
  const [position, setPosition] = useState<number | null>(null)
  const [waitTime, setWaitTime] = useState<number | null>(null)

  useEffect(() => {
    const joinQueue = async () => {
      // const res = await fetch("/api/queue",{
      //   method : "POST"
      // });
      // const data = await res.json();
      // setTicketId(data.ticketId)
      // setPosition(5)
      // setWaitTime(100)
       const fakeData = {
        ticketId: "dummy-" + Math.random().toString(36).slice(2, 8),
        position: 5,
        waitTime: 100,
      }
    }
    joinQueue();
  },[])

  useEffect(() => {
    if(!ticketId) return

    const interval = setInterval(async() => {
      // const res = await fetch(`/api/queue?ticketId=${ticketId}`)
      // const data = await res.json()
      // setPosition(data.position)
      // setWaitTime(data.waitTime)
      setPosition((prev) => prev !== null ? Math.max(prev - 1, 0) : null)
      setWaitTime((prev) => prev !== null ? Math.max(prev - 1, 0) : null)
    },3000)
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
