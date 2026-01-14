import { useRef, useState } from 'react'

function FileUpload({ label, icon, file, onFileSelect, subtitle }) {
    const inputRef = useRef(null)
    const [isDragOver, setIsDragOver] = useState(false)

    const handleDragOver = (e) => {
        e.preventDefault()
        setIsDragOver(true)
    }

    const handleDragLeave = () => {
        setIsDragOver(false)
    }

    const handleDrop = (e) => {
        e.preventDefault()
        setIsDragOver(false)
        const droppedFile = e.dataTransfer.files[0]
        if (droppedFile?.name.endsWith('.csv')) {
            onFileSelect(droppedFile)
        }
    }

    const handleClick = () => {
        inputRef.current?.click()
    }

    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0]
        if (selectedFile) {
            onFileSelect(selectedFile)
        }
    }

    return (
        <div
            className={`upload-zone ${isDragOver ? 'dragover' : ''} ${file ? 'has-file' : ''}`}
            onClick={handleClick}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
        >
            <input
                ref={inputRef}
                type="file"
                accept=".csv"
                onChange={handleFileChange}
                style={{ display: 'none' }}
            />
            <div className="upload-icon">{file ? '✅' : icon}</div>
            <div className="upload-title">{label}</div>
            <div className="upload-subtitle">
                {file ? '' : subtitle}
            </div>
            {file && <div className="file-name">{file.name}</div>}
            {!file && (
                <div className="upload-subtitle" style={{ marginTop: '12px' }}>
                    Drag & drop or click to browse
                </div>
            )}
        </div>
    )
}

export default FileUpload
