'use client';

import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';

interface PacketSize {
  sizes: number[];
}

interface CalculateResponse {
  data: {
    optimal_packets: Record<string, number>;
  };
}

export default function Home() {
  const [packetSizes, setPacketSizes] = useState<number[]>([]);
  const [newSizes, setNewSizes] = useState<string>('');
  const [items, setItems] = useState<string>('');
  const [optimalPackets, setOptimalPackets] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [retryCount, setRetryCount] = useState<number>(0);
  const [isConnected, setIsConnected] = useState<boolean>(false);

  useEffect(() => {
    fetchPacketSizes();
  }, []);

  const fetchPacketSizes = async () => {
    try {
      setLoading(true);
      
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout
      
      const response = await fetch('/api/v1/packet/size', {
        signal: controller.signal,
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      
      const data = await response.json();
      setPacketSizes(data.data.packet_sizes);
      setIsConnected(true);
      setError('');
    } catch (err) {
      console.error('Error fetching packet sizes:', err);
      
      if (err instanceof Error && err.name === 'AbortError') {
        setError('Request timed out. The server took too long to respond.');
      } else {
        setError('Failed to connect to the server. Please make sure the backend is running on port 3000.');
      }
      
      setIsConnected(false);
      
      if (retryCount < 3) {
        const backoffDelay = Math.min(1000 * Math.pow(2, retryCount), 10000); // Max 10s
        setTimeout(() => {
          setRetryCount(prev => prev + 1);
          fetchPacketSizes();
        }, backoffDelay);
      }
    } finally {
      setLoading(false);
    }
  };

  const updatePacketSizes = async () => {
    try {
      setLoading(true);
      const sizes = newSizes.split(',').map(size => parseInt(size.trim())).filter(size => !isNaN(size));
      
      if (sizes.length === 0) {
        setError('Please enter valid numbers separated by commas');
        return;
      }
      
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 15000); // 15 second timeout
      
      const response = await fetch('/api/v1/packet/size', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ sizes }),
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      
      await fetchPacketSizes();
      setNewSizes('');
      setError('');
    } catch (err) {
      console.error('Error updating packet sizes:', err);
      
      if (err instanceof Error && err.name === 'AbortError') {
        setError('Request timed out. The server took too long to respond.');
      } else {
        setError('Failed to update packet sizes. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  const calculateOptimalPackets = async () => {
    try {
      if (!items || isNaN(parseInt(items)) || parseInt(items) <= 0) {
        setError('Please enter a valid number of items');
        return;
      }
      
      const itemCount = parseInt(items);
      const maxAllowed = 1_000_000_000; // 1 billion
      
      if (itemCount > maxAllowed) {
        setError(`The maximum allowed number of items is ${maxAllowed.toLocaleString()} (1 billion). Please enter a smaller number.`);
        return;
      }
      
      // if (itemCount > 10_000_000) {
      //   const proceed = window.confirm(
      //     `You entered ${itemCount.toLocaleString()} items. Large calculations may take 10-30 seconds. Do you want to proceed?`
      //   );
      //   if (!proceed) {
      //     return;
      //   }
      // }
      
      setLoading(true);
      setError('');
      
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 60000); // 60 second timeout for calculations
      
      const response = await fetch(`/api/v1/packet/calculate?items=${items}`, {
        signal: controller.signal,
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error! Status: ${response.status} - ${errorText}`);
      }
      
      const data: CalculateResponse = await response.json();
      setOptimalPackets(data.data.optimal_packets);
      setError('');
    } catch (err) {
      console.error('Error calculating optimal packets:', err);
      
      if (err instanceof Error && err.name === 'AbortError') {
        setError('Calculation timed out after 60 seconds. Try a smaller number of items or check if the backend is overloaded.');
      } else if (err instanceof Error) {
        setError(`Failed to calculate optimal packets: ${err.message}`);
      } else {
        setError('Failed to calculate optimal packets. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-900 text-white p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <motion.div 
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center mb-12"
        >
          <h1 className="text-5xl font-bold mb-4 bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-purple-400">
            Packet Size Manager
          </h1>
          <p className="text-gray-300 text-lg">
            Manage and optimize your packet sizes with ease
          </p>
        </motion.div>

        {!isConnected && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="bg-red-500/20 border border-red-500 text-white p-4 rounded-lg flex items-center justify-between"
          >
            <div className="flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
              <span>Connection to backend server failed. Retrying... ({retryCount}/3)</span>
            </div>
            <button 
              onClick={fetchPacketSizes}
              className="bg-red-500 hover:bg-red-600 px-4 py-1 rounded-md text-sm transition-colors"
            >
              Retry Now
            </button>
          </motion.div>
        )}

        {error && isConnected && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="bg-red-500/20 border border-red-500 text-white p-4 rounded-lg"
          >
            {error}
          </motion.div>
        )}

        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white/10 backdrop-blur-md p-6 rounded-xl shadow-lg border border-white/20"
        >
          <h2 className="text-2xl font-semibold mb-4 flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            Current Packet Sizes
          </h2>
          <div className="flex flex-wrap gap-2">
            {packetSizes.length > 0 ? (
              packetSizes.map((size) => (
                <span key={size} className="bg-blue-500/30 border border-blue-400 px-4 py-2 rounded-full text-blue-100">
                  {size}
                </span>
              ))
            ) : (
              <p className="text-gray-400 italic">No packet sizes available</p>
            )}
          </div>
        </motion.div>

        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white/10 backdrop-blur-md p-6 rounded-xl shadow-lg border border-white/20"
        >
          <h2 className="text-2xl font-semibold mb-4 flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Update Packet Sizes
          </h2>
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1">
              <label htmlFor="newSizes" className="block text-sm font-medium text-gray-300 mb-1">
                Enter sizes (comma-separated)
              </label>
              <input
                id="newSizes"
                type="text"
                value={newSizes}
                onChange={(e) => setNewSizes(e.target.value)}
                placeholder="e.g., 100, 200, 500"
                className="w-full bg-gray-800/50 border border-gray-700 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-purple-500 text-white"
              />
            </div>
            <div className="flex items-end">
              <button
                onClick={updatePacketSizes}
                disabled={loading || !isConnected}
                className="bg-purple-600 hover:bg-purple-700 px-6 py-2 rounded-lg transition-colors disabled:opacity-50 flex items-center"
              >
                {loading ? (
                  <>
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Updating...
                  </>
                ) : (
                  'Update'
                )}
              </button>
            </div>
          </div>
        </motion.div>

        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white/10 backdrop-blur-md p-6 rounded-xl shadow-lg border border-white/20"
        >
          <h2 className="text-2xl font-semibold mb-4 flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
            Calculate Optimal Packets
          </h2>
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1">
              <label htmlFor="items" className="block text-sm font-medium text-gray-300 mb-1">
                Number of items (max: 1 billion)
              </label>
              <input
                id="items"
                type="number"
                value={items}
                onChange={(e) => setItems(e.target.value)}
                placeholder="e.g., 10000000"
                min="1"
                max="1000000000"
                className="w-full bg-gray-800/50 border border-gray-700 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-green-500 text-white"
              />
              <p className="text-xs text-gray-400 mt-1">
                💡 Tip: Calculations over 10 million items may take 10-30 seconds
              </p>
            </div>
            <div className="flex items-end">
              <button
                onClick={calculateOptimalPackets}
                disabled={loading || !isConnected}
                className="bg-green-600 hover:bg-green-700 px-6 py-2 rounded-lg transition-colors disabled:opacity-50 flex items-center"
              >
                {loading ? (
                  <>
                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Calculating...
                  </>
                ) : (
                  'Calculate'
                )}
              </button>
            </div>
          </div>
          
          {Object.keys(optimalPackets).length > 0 && (
            <div className="mt-6">
              <h3 className="text-xl font-semibold mb-3 flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Results
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {Object.entries(optimalPackets).map(([size, count]) => (
                  <div key={size} className="bg-gray-800/50 p-4 rounded-lg border border-gray-700">
                    <div className="flex justify-between items-center">
                      <div>
                        <p className="text-sm text-gray-400">Size</p>
                        <p className="text-lg font-medium">{size}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-gray-400">Count</p>
                        <p className="text-lg font-medium">{count}</p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </motion.div>
      </div>
    </main>
  );
} 