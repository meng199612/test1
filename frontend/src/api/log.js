import request from '@/utils/request'

export function getOperationLogs(params) {
  return request.get('/logs', { params })
}

export function getOperationLog(id) {
  return request.get(`/logs/${id}`)
}

export function deleteOperationLog(id) {
  return request.delete(`/logs/${id}`)
}

export function clearOperationLogs(days) {
  return request.post('/logs/clear', { days })
}
