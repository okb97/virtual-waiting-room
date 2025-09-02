"use client"

import {useEffect, useState} from 'react';

export default function Home() {
  const [ticketId, setTicketId] = useState<string | null>(null)

  useEffect(() => {
    const joinQueue = async () => {
      const res = await fetch("/api/queue",{
        method : "POST"
      });
      const data = await res.json();
      setTicketId(data.ticketId)
    }
    joinQueue();
  },[])


  return (
    <main>
      <h1>仮装待合室</h1>
      <p>あなたのチケットID:{ticketId}</p>
    </main>
  );
}
