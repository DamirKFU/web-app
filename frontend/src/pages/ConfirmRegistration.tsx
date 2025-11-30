import { useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { addCSRFToken } from "@/lib/csrf";
import { API_CONFIG } from "@/config/api";
import { useToast } from "@/hooks/use-toast";

const API_BASE_URL = API_CONFIG.baseURL;

export default function ConfirmRegistration() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: "Ссылка недействительна. Пожалуйста, зарегистрируйтесь снова.",
      });
      navigate("/login", { replace: true });
      return;
    }

    const confirm = async () => {
      try {
        const headers = await addCSRFToken();

        const response = await fetch(`${API_BASE_URL}/api/v1/auth/register-confirm/`, {
          method: "POST",
          headers: {
            ...headers,
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ token }),
        });

        if (response.ok) {
          toast({
            title: "Успешно",
            description: "Ваш email успешно подтверждён. Теперь вы можете войти.",
          });
        } else {
          toast({
            variant: "destructive",
            title: "Ошибка подтверждения",
            description: "Ссылка недействительна или устарела. Пожалуйста, зарегистрируйтесь снова.",
          });
        }

      } catch {
        toast({
          variant: "destructive",
          title: "Ошибка",
          description: "Не удалось подтвердить регистрацию. Повторите регистрацию.",
        });
      } finally {
        navigate("/login", { replace: true });
      }
    };

    confirm();
  }, [searchParams, navigate, toast]);

  return (
    <div className="flex items-center justify-center h-[60vh] text-lg">
      Подтверждение регистрации...
    </div>
  );
}
