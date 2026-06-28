import { useQuery } from "@tanstack/react-query";
import { getAuditLogs } from "../api/audit_logs";

export const useGetAuditLogs = () => {
    return useQuery({
        queryKey: ["audits"],
        queryFn: ()=> getAuditLogs(),
    })
}