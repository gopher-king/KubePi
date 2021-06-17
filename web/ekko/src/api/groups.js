import {post, get} from "@/plugins/request"

const baseUrl = "/api/v1/groups"

export function searchGroups(pageNum, pageSize, conditions) {
    return post(`${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`, conditions)
}

export function listGroups() {
    return get(`${baseUrl}`)
}