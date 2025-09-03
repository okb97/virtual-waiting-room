"use client"

import {useEffect, useState} from 'react';

export default function Home() {
  const [ticketId, setTicketId] = useState<string | null>(null)
  const [position, setPosition] = useState<number | null>(null)
  const [waitTime, setWaitTime] = useState<number | null>(null)

  useEffect(() => {
    const joinQueue = async () => {
      const res = await fetch("/api/queue",{
        method : "POST"
      });
      const data = await res.json();
      setTicketId(data.ticketId)
      setPosition(5)
      setWaitTime(100)
    }
    joinQueue();
  },[])

  useEffect(() => {
    if(!ticketId) return

    const interval = setInterval(() => {
      setPosition((prev) => (prev && prev > 1 ? prev - 1 : 0))
      setWaitTime((prev) => (prev && prev > 0 ? prev - 1 : 0))
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
