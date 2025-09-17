"use client"

import { useRouter } from 'next/navigation';
import {useEffect, useState} from 'react';

export default function Home() {
  const router = useRouter()
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
    },10000)
    return () => clearInterval(interval)
  },[ticketId])

  useEffect(() => {
    if (!ticketId) return;
    const interval = setInterval(async () => {
      const res = await fetch(`/api/eligible?ticketId=${ticketId}`);
      const data = await res.json();
      if (data.canPurchase) {
        clearInterval(interval);
        router.push('/purchase');
      }
    }, 30000);
    return () => clearInterval(interval);
  }, [ticketId, router]);

  return (
    <main>
      <h1>仮想待合室</h1>
      <p className="mb-2">あなたのチケットID: <strong>{ticketId}</strong></p>
          <p className="mb-2">現在の順番: {position !== null ? position + 1 : "順番探索中"} 番目</p>
          <p className="mb-4">推定待ち時間: {waitTime !== null ? waitTime : "待ち時間探索中"} 分</p>
    </main>
  );
}
