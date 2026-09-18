'use client';
import { useState, useEffect } from 'react';
import {
  LanguageIcon,
  ArrowsRightLeftIcon,
  ClipboardIcon,
  ClockIcon,
  StarIcon,
  TrashIcon,
  SparklesIcon,
} from '@heroicons/react/24/outline';

interface Language {
  code: string;
  name: string;
}

interface TranslationRecord {
  id: string;
  text: string;
  source_lang: string;
  target_lang: string;
  translated: string;
  created_at: string;
}

interface Term {
  id: string;
  source: string;
  target: string;
  lang_pair: string;
}

export default function Home() {
  const [text, setText] = useState('');
  const [sourceLang, setSourceLang] = useState('auto');
  const [targetLang, setTargetLang] = useState('zh');
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);
  const [languages, setLanguages] = useState<Language[]>([]);
  const [history, setHistory] = useState<TranslationRecord[]>([]);
  const [terms, setTerms] = useState<Term[]>([]);
  const [activeTab, setActiveTab] = useState('translate');
  const [swapLoading, setSwapLoading] = useState(false);

  useEffect(() => {
    fetch('/api/languages')
      .then(r => r.json())
      .then(setLanguages)
      .catch(console.error);

    fetch('/api/history?limit=20')
      .then(r => r.json())
      .then(setHistory)
      .catch(console.error);

    fetch('/api/terms')
      .then(r => r.json())
      .then(setTerms)
      .catch(console.error);
  }, []);

  const handleTranslate = async () => {
    if (!text.trim()) return;
    
    setLoading(true);
    try {
      const res = await fetch('/api/translate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          text,
          source: sourceLang,
          target: targetLang,
        }),
      });
      const data = await res.json();
      setResult(data.translated_text || data.translated || '');
    } catch (err) {
      console.error('Translation failed:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSwap = () => {
    setSwapLoading(true);
    setSourceLang(targetLang);
    setTargetLang(sourceLang);
    setResult(text);
    setText(result);
    setTimeout(() => setSwapLoading(false), 300);
  };

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-blue-900 to-slate-900">
      <div className="max-w-6xl mx-auto px-4 py-12">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold text-white mb-4 flex items-center justify-center gap-3">
            <LanguageIcon className="w-12 h-12 text-blue-400" />
            AI Translator
          </h1>
          <p className="text-xl text-blue-200">智能翻译 · 多语言互译 · 术语管理</p>
        </div>

        {/* Tabs */}
        <div className="flex justify-center gap-4 mb-8">
          {['translate', 'history', 'terms'].map(tab => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-6 py-2 rounded-full font-medium transition-all ${
                activeTab === tab 
                  ? 'bg-blue-600 text-white' 
                  : 'bg-slate-800/50 text-gray-400 hover:text-white'
              }`}
            >
              {tab === 'translate' && '翻译'}
              {tab === 'history' && '历史'}
              {tab === 'terms' && '术语库'}
            </button>
          ))}
        </div>

        {/* Translate Tab */}
        {activeTab === 'translate' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Source */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-blue-500/20">
              <div className="flex justify-between items-center mb-4">
                <select
                  value={sourceLang}
                  onChange={e => setSourceLang(e.target.value)}
                  className="px-4 py-2 bg-slate-900/50 border border-blue-500/30 rounded-lg text-white focus:outline-none focus:border-blue-500"
                >
                  {languages.map(lang => (
                    <option key={lang.code} value={lang.code}>
                      {lang.name}
                    </option>
                  ))}
                </select>
                <button
                  onClick={() => handleCopy(text)}
                  className="p-2 text-gray-400 hover:text-white transition-colors"
                >
                  <ClipboardIcon className="w-5 h-5" />
                </button>
              </div>
              <textarea
                value={text}
                onChange={e => setText(e.target.value)}
                placeholder="输入要翻译的文本..."
                className="w-full h-64 px-4 py-3 bg-slate-900/50 border border-blue-500/30 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:border-blue-500 resize-none"
              />
              <div className="text-right text-gray-400 text-sm mt-2">
                {text.length} 字符
              </div>
            </div>

            {/* Target */}
            <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-blue-500/20">
              <div className="flex justify-between items-center mb-4">
                <select
                  value={targetLang}
                  onChange={e => setTargetLang(e.target.value)}
                  className="px-4 py-2 bg-slate-900/50 border border-blue-500/30 rounded-lg text-white focus:outline-none focus:border-blue-500"
                >
                  {languages.filter(l => l.code !== 'auto').map(lang => (
                    <option key={lang.code} value={lang.code}>
                      {lang.name}
                    </option>
                  ))}
                </select>
                <button
                  onClick={() => handleCopy(result)}
                  className="p-2 text-gray-400 hover:text-white transition-colors"
                >
                  <ClipboardIcon className="w-5 h-5" />
                </button>
              </div>
              <div className="w-full h-64 px-4 py-3 bg-slate-900/50 border border-blue-500/30 rounded-xl text-white overflow-y-auto">
                {loading ? (
                  <div className="flex items-center justify-center h-full">
                    <SparklesIcon className="w-8 h-8 text-blue-400 animate-spin" />
                  </div>
                ) : result ? (
                  result
                ) : (
                  <span className="text-gray-500">翻译结果将显示在这里...</span>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Swap Button */}
        <div className="flex justify-center my-6">
          <button
            onClick={handleSwap}
            disabled={swapLoading || !text}
            className="p-3 bg-slate-700 hover:bg-slate-600 disabled:bg-slate-800 disabled:text-gray-500 text-white rounded-full transition-all"
          >
            <ArrowsRightLeftIcon className={`w-6 h-6 ${swapLoading ? 'animate-spin' : ''}`} />
          </button>
        </div>

        {/* Translate Button */}
        <div className="flex justify-center mb-12">
          <button
            onClick={handleTranslate}
            disabled={loading || !text.trim()}
            className="px-12 py-4 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 disabled:from-slate-600 disabled:to-slate-700 text-white rounded-xl font-semibold text-lg transition-all flex items-center gap-2"
          >
            <SparklesIcon className="w-6 h-6" />
            {loading ? '翻译中...' : '开始翻译'}
          </button>
        </div>

        {/* History Tab */}
        {activeTab === 'history' && (
          <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-blue-500/20">
            <h2 className="text-xl font-semibold text-white mb-4 flex items-center gap-2">
              <ClockIcon className="w-6 h-6 text-blue-400" />
              翻译历史
            </h2>
            {history.length === 0 ? (
              <p className="text-gray-400 text-center py-8">暂无翻译历史</p>
            ) : (
              <div className="space-y-4">
                {history.map(record => (
                  <div key={record.id} className="p-4 bg-slate-900/50 rounded-xl">
                    <div className="flex justify-between items-start mb-2">
                      <span className="text-xs text-gray-500">
                        {record.source_lang} → {record.target_lang}
                      </span>
                      <span className="text-xs text-gray-500">
                        {new Date(record.created_at).toLocaleString()}
                      </span>
                    </div>
                    <p className="text-gray-300 mb-2">{record.text}</p>
                    <p className="text-blue-300">{record.translated}</p>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Terms Tab */}
        {activeTab === 'terms' && (
          <div className="bg-slate-800/50 backdrop-blur rounded-2xl p-6 border border-blue-500/20">
            <h2 className="text-xl font-semibold text-white mb-4 flex items-center gap-2">
              <StarIcon className="w-6 h-6 text-yellow-400" />
              术语库
            </h2>
            {terms.length === 0 ? (
              <p className="text-gray-400 text-center py-8">暂无术语</p>
            ) : (
              <div className="space-y-3">
                {terms.map(term => (
                  <div key={term.id} className="p-4 bg-slate-900/50 rounded-xl flex justify-between items-center">
                    <div>
                      <span className="text-gray-400 text-sm">{term.lang_pair}</span>
                      <p className="text-white">{term.source} → {term.target}</p>
                    </div>
                    <button
                      onClick={() => {
                        fetch(`/api/terms/${term.id}`, { method: 'DELETE' });
                        setTerms(terms.filter(t => t.id !== term.id));
                      }}
                      className="p-2 text-red-400 hover:text-red-300"
                    >
                      <TrashIcon className="w-5 h-5" />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
